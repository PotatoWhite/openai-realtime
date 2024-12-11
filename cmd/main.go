package main

import (
	"bufio"
	"context"
	"github.com/gordonklaus/portaudio"
	"github.com/sirupsen/logrus"
	"openai-realtime/pkg/audiomanager"
	"openai-realtime/pkg/audioutils"
	"openai-realtime/pkg/config"
	"openai-realtime/pkg/openai"
	"openai-realtime/pkg/openai/events"
	"os"
	"os/signal"
	"time"
)

var (
	log = func() *logrus.Logger {
		log := logrus.New()
		log.SetOutput(os.Stdout)
		log.SetLevel(config.LogLevel)
		return log
	}()

	openAIApiKey          = os.Getenv("OPENAI_API_KEY")
	openAIAudioSampleRate = 24000
	txFileName            = "tx.pcm"
	rxFileName            = "rx.pcm"
	txWavFileName         = "tx.wav"
	rxWavFileName         = "rx.wav"
)

func initializePortAudio() {
	if err := portaudio.Initialize(); err != nil {
		log.Fatalf("Failed to initialize PortAudio: %v", err)
	}
}

func shutdownPortAudio() {
	if err := portaudio.Terminate(); err != nil {
		log.Errorf("Failed to terminate PortAudio: %v", err)
	}
}

func createOpenAIClient(ctx context.Context) *openai.Client {
	log.Info("Creating OpenAI client")
	client, err := openai.NewClient(ctx, "api.openai.com", "/v1/realtime", "gpt-4o-realtime-preview-2024-10-01", openAIApiKey)
	if err != nil {
		log.Fatalf("Failed to create OpenAI client: %v", err)
	}
	log.Info("OpenAI client created successfully")
	return client
}

// Enter 키를 누르면 종료 신호를 보내는 함수
func waitForUserExitSignal(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		log.Debug("Wait for user exit signal stopped")
	}()

	reader := bufio.NewReader(os.Stdin)
	inputCh := make(chan string)
	errCh := make(chan error)

	// 사용자 입력을 읽는 별도의 goroutine
	go func() {
		input, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		inputCh <- input
	}()

	select {
	case <-ctx.Done():
		// 컨텍스트가 취소되면 함수 종료
		return
	case input := <-inputCh:
		log.Info("Enter key pressed")
		// 필요에 따라 입력값을 처리할 수 있습니다.
		_ = input
		cancel()
		log.Info("Context cancelled by user key press")
	case err := <-errCh:
		log.Errorf("Error reading input: %v", err)
	}
}

func handleInterruptSignal(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		log.Debug("Handle interrupt signal stopped")
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-ctx.Done():
		// 컨텍스트가 취소되면 함수 종료
		return
	case <-interrupt:
		log.Info("Interrupt signal received")
		cancel()
	}
}

func listenAndSendToOpenAI(ctx context.Context, am *audiomanager.Manager, openAI *openai.Client, cancel context.CancelFunc) {
	defer func() {
		log.Debug("Audio processing to OpenAI stopped")
	}()
	log.Info("Starting audio processing to OpenAI")

	for {
		select {
		case <-ctx.Done():
			log.Info("Context done, stopping audio processing")
			am.Close()
			return
		case audioData, ok := <-am.DeviceController.InputChan:
			if !ok {
				log.Info("Audio channel closed, stopping audio processing")
				return
			}

			// Ignore audio below the silence threshold
			rmsThresholdDb := config.RmsThresholdDb
			useZCR := config.UseZCR
			zcrThreshold := config.ZcrThreshold

			if audiomanager.IsSilentAudioDataEx(audioData, rmsThresholdDb, useZCR, zcrThreshold) {
				log.Debug("Silent audio data detected, skipping transmission")
				continue
			}

			// Resample audio to 24000 Hz for OpenAI
			resampled := audioutils.ResampleAudioData(audioData, am.DeviceController.SampleRate, openAIAudioSampleRate)
			if resampled == nil {
				log.Error("Resampling failed, skipping this chunk")
				continue
			}
			log.Debug("Resampled audio data from %d Hz to 24000 Hz", am.DeviceController.SampleRate)

			// Convert to byte array
			byteAudioData := audioutils.ConvertToByteArrayLE(resampled)
			log.Debug("Converted resampled audio data to byte array")

			// Send audio data to OpenAI
			if err := openAI.SendInputAudioBufferAppend(byteAudioData); err != nil {
				log.Errorf("Failed to send audio to OpenAI: %v", err)
				cancel() // Cancel context on error
				return
			}

			// Append PCM data to file in a separate goroutine
			go func(data []byte) {
				if err := openai.AppendPCMDataToFile(txFileName, data); err != nil {
					log.Errorf("Failed to append PCM data to file: %v", err)
					cancel() // Cancel context on error
				}
			}(byteAudioData)
		}
	}
}
func receiveAndSaveFromOpenAI(ctx context.Context, am *audiomanager.Manager, openAI *openai.Client, cancel context.CancelFunc) {
	defer func() {
		log.Debug("Receive and save from OpenAI stopped")
	}()
	log.Info("Starting to receive and save from OpenAI")

	framesPerBuffer := am.DeviceController.SampleRate / 10 // e.g., 44100 / 10 = 4410
	buffer := make([]int16, 0, framesPerBuffer*10)         // Initialize buffer with sufficient capacity

	for {
		select {
		case <-ctx.Done():
			log.Info("Context done, stopping receive and save from OpenAI")
			return
		case audioData, ok := <-openAI.AudioOutputChan:
			if !ok {
				log.Info("Audio channel closed, stopping receive and save from OpenAI")
				return
			}

			log.Debugf("Received audio data of length: %d bytes", len(audioData))

			// Convert byte array to int16 PCM samples
			pcmSamples := audioutils.ConvertToInt16ArrayLE(audioData)
			if pcmSamples == nil {
				log.Error("PCM conversion failed, skipping this chunk")
				continue
			}

			// Resample audio to match Manager's sample rate
			resampled := audioutils.ResampleAudioData(pcmSamples, openAIAudioSampleRate, am.DeviceController.SampleRate)
			if resampled == nil {
				log.Error("Resampling failed, skipping this chunk")
				continue
			}
			log.Debugf("Resampled audio length: %d samples", len(resampled))

			// Append resampled data to the buffer
			buffer = append(buffer, resampled...)

			// Send fixed-size chunks to OutputChan for playback
			for len(buffer) >= framesPerBuffer {
				chunk := buffer[:framesPerBuffer]
				buffer = buffer[framesPerBuffer:]

				select {
				case <-ctx.Done():
					log.Info("Context done while sending to OutputChan")
					return
				case am.DeviceController.OutputChan <- chunk:
					log.Debug("Sending buffer of length: %d samples to OutputChan", len(chunk))
				}
			}

			// Append PCM data to file in a separate goroutine
			go func(data []byte) {
				if err := openai.AppendPCMDataToFile(rxFileName, data); err != nil {
					log.Errorf("Failed to append PCM data to file: %v", err)
					cancel() // Optionally cancel context on error
				}
			}(audioData)
		}
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 오디오장치 초기화
	initializePortAudio()
	defer shutdownPortAudio()

	inputDevice := audiomanager.SelectInputDevice()   // 입력 장치 선택
	outputDevice := audiomanager.SelectOutputDevice() // 출력 장치 선택

	// 오디오 매니저 생성
	audioManager, err := audiomanager.NewManager(inputDevice, outputDevice, 10)
	if err != nil {
		log.Fatalf("Failed to create audio manager: %v", err)
	}
	defer audioManager.Close()

	// OpenAI 클라이언트 생성
	openAI := createOpenAIClient(ctx)
	defer openAI.Close()

	// OpenAI 에 Project 전송
	iat := events.InputAudioTranscription{
		Model: "whisper-1",
	}

	tDetection := events.TurnDetection{
		Type:              "server_vad",
		Threshold:         0.5,
		PrefixPaddingMs:   300,
		SilenceDurationMs: 500,
	}

	openAI.SessionUpdate(iat, tDetection, []events.Tool{})

	// 파일 명 업데이트 (날짜_파일명)
	datetime := time.Now().Format("20060102_150405")
	txFileName = datetime + "_" + txFileName
	rxFileName = datetime + "_" + rxFileName
	txWavFileName = datetime + "_" + txWavFileName
	rxWavFileName = datetime + "_" + rxWavFileName

	// ReceiveServerEvent goroutine
	go openAI.ReceiveServerEvent(ctx, cancel)                      // openAI의 ServerEvent 를 수신 및 처리
	go audioManager.Start(ctx)                                     // 오디오 매니저 시작
	go listenAndSendToOpenAI(ctx, audioManager, openAI, cancel)    // 오디오 장치로부터 오디오를 받아 OpenAI로 전송
	go receiveAndSaveFromOpenAI(ctx, audioManager, openAI, cancel) // OpenAI로부터 오디오를 받아 재생

	go waitForUserExitSignal(ctx, cancel) // 사용자 입력을 대기 및 종료 신호 전달
	go handleInterruptSignal(ctx, cancel) // 인터럽트 신호 수신 및 종료 신호 전달

	// 종료 신호를 대기
	<-ctx.Done()
	log.Info("Context done in main function")

	// convert pcm to wav by ffmpeg

	if err := audioutils.ConvertPCMToWav(rxFileName, rxWavFileName, 24000, 1); err != nil {
		log.Errorf("Failed to convert PCM to WAV: %v", err)
	}

	if err := audioutils.ConvertPCMToWav(txFileName, txWavFileName, 24000, 1); err != nil {
		log.Errorf("Failed to convert PCM to WAV: %v", err)
	}

	//wg.Wait()
	log.Info("Program terminated")
}
