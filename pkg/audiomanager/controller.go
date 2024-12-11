package audiomanager

import (
	"context"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"openai-realtime/pkg/audioutils"
	"sync"
)

// Recorder는 오디오를 녹음하는 구조체입니다.
type Controller struct {
	InputDevice  *portaudio.DeviceInfo
	OutputDevice *portaudio.DeviceInfo
	SampleRate   int
	VolumeThresh float32

	InputChan  chan []int16 // 마이크로부터 오디오 데이터 채널
	OutputChan chan []int16 // 스피커로 출력할 오디오 데이터 채널
	ErrorChan  chan error   // 오류 채널

	stream   *portaudio.Stream // 포트오디오 스트림
	stopOnce sync.Once         // Off() 메서드가 한 번만 실행되도록 보장
}

// NewController 생성자 함수
func NewController(inputDevice *portaudio.DeviceInfo, outputDevice *portaudio.DeviceInfo, volumeThreshold float32) *Controller {
	return &Controller{
		InputDevice:  inputDevice,
		OutputDevice: outputDevice,
		SampleRate:   int(inputDevice.DefaultSampleRate),
		VolumeThresh: volumeThreshold,
		InputChan:    make(chan []int16, 5), // 버퍼링하여 블로킹 방지
		OutputChan:   make(chan []int16, 5), // 버퍼링하여 블로킹 방지
		ErrorChan:    make(chan error, 1),
	}
}

// Off 녹음을 중지하고 스트림을 종료합니다.
func (c *Controller) Off() {
	c.stopOnce.Do(func() {
		if c.stream != nil {
			if err := c.stream.Abort(); err != nil {
				log.Errorf("Failed to abort stream: %v", err)
			}
		}
		close(c.ErrorChan) // 에러 채널을 닫아 더 이상의 에러 전송을 방지
		log.Info("Controller.Off")
	})
}

// RecordAudio는 오디오를 녹음하는 함수입니다.
func (c *Controller) On(ctx context.Context) error {
	defer log.Debug("Controller stopped")

	// 스트림 열기
	stream, err := portaudio.OpenStream(c.getStreamParam(), c.process)
	if err != nil {
		close(c.InputChan)
		return fmt.Errorf("failed to open stream: %w", err)
	}
	c.stream = stream
	defer func() {
		if err := stream.Close(); err != nil {
			log.Warnf("Failed to close stream: %v", err)
		}
	}()

	if err := stream.Start(); err != nil {
		close(c.InputChan)
		return fmt.Errorf("failed to start stream: %w", err)
	}
	defer func() {
		if err := stream.Stop(); err != nil {
			log.Warnf("Failed to stop stream: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info("Controller received context cancellation before read")
			c.Off()
			close(c.InputChan)
			return nil
		case err := <-c.ErrorChan:
			log.Warn("Controller received error before read")
			c.Off()
			close(c.InputChan)
			return err
		}
	}
}

func (c *Controller) getStreamParam() portaudio.StreamParameters {
	if c.InputDevice == nil {
		log.Fatal("Input device is not set")
	}

	if c.OutputDevice == nil {
		log.Fatal("Output device is not set")
	}

	return portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   c.InputDevice,
			Channels: 1,
		},
		Output: portaudio.StreamDeviceParameters{
			Device:   c.OutputDevice,
			Channels: 1, // mono 출력
			Latency:  c.OutputDevice.DefaultLowOutputLatency,
		},
		SampleRate:      float64(c.SampleRate),
		FramesPerBuffer: c.SampleRate / 10, // 0.1초 단위 버퍼
	}
}

func (c *Controller) process(in []int16, out []int16) {
	// 입력 처리
	inputCopy := make([]int16, len(in))
	copy(inputCopy, in)
	select {
	case c.InputChan <- inputCopy:
	default:
		log.Warn("Input channel is full, discarding audio data")
	}

	// 출력 처리
	select {
	case data := <-c.OutputChan:
		audioutils.CopyAudioData(out, data)
	default:
		for i := range out {
			out[i] = 0
		}
	}
}
