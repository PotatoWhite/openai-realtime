package audioutils

import (
	"fmt"
	"os/exec"
)

// ConvertToInt16Array converts byte array to int16 array (little-endian)
func ConvertToInt16ArrayLE(data []byte) []int16 {
	audioData := make([]int16, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		audioData[i/2] = int16(data[i]) | int16(data[i+1])<<8
	}
	return audioData
}

// ConvertToInt16ArrayBE converts byte array to int16 array (big-endian)
func ConvertToInt16ArrayBE(data []byte) []int16 {
	audioData := make([]int16, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		audioData[i/2] = int16(data[i+1]) | int16(data[i])<<8
	}
	return audioData
}

// @title resampleAudioData
// @description 오디오 데이터를 재샘플링하는 함수
// @param data []int16 오디오 데이터
// @param fromRate int 원본 샘플 속도
// @param toRate int 대상 샘플 속도
// @return []int16 재샘플링된 오디오 데이터
func ResampleAudioData(data []int16, fromRate int, toRate int) []int16 {
	ratio := float64(toRate) / float64(fromRate)
	newLength := int(float64(len(data)) * ratio)
	resampledData := make([]int16, newLength)

	for i := 0; i < newLength; i++ {
		srcIndex := float64(i) / ratio
		srcIndexInt := int(srcIndex)
		nextIndex := srcIndexInt + 1
		srcIndexFrac := srcIndex - float64(srcIndexInt)

		// Linear interpolation with boundary check
		if nextIndex < len(data) {
			resampledData[i] = int16(float64(data[srcIndexInt])*(1-srcIndexFrac) +
				float64(data[nextIndex])*srcIndexFrac)
		} else {
			resampledData[i] = data[srcIndexInt]
		}
	}
	return resampledData
}

// copyAudioData는 데이터를 출력 버퍼로 복사합니다.
func CopyAudioData(out, data []int16) {
	if len(data) > len(out) {
		data = data[:len(out)]
	}
	copy(out, data)
	if len(data) < len(out) {
		fillSilence(out[len(data):]) // 데이터 길이가 부족한 경우 무음으로 채움
	}
}

// fillSilence는 버퍼를 무음으로 채웁니다.
func fillSilence(buffer []int16) {
	for i := range buffer {
		buffer[i] = 0
	}
}

// Start starts the audio manager
func ChangeEndian(data []byte) []byte {
	for i := 0; i < len(data); i += 2 {
		data[i], data[i+1] = data[i+1], data[i]
	}
	return data
}

// ConvertToByteArrayLE converts int16 array to byte array (little-endian)
func ConvertToByteArrayLE(data []int16) []byte {
	audioBytes := make([]byte, len(data)*2)
	for i, sample := range data {
		audioBytes[i*2] = byte(sample & 0xFF)
		audioBytes[i*2+1] = byte(sample >> 8)
	}
	return audioBytes
}

// ConvertToByteArrayBE converts int16 array to byte array (big-endian)
func ConvertToByteArrayBE(data []int16) []byte {
	audioBytes := make([]byte, len(data)*2)
	for i, sample := range data {
		audioBytes[i*2] = byte(sample >> 8)
		audioBytes[i*2+1] = byte(sample & 0xFF)
	}
	return audioBytes
}

// ConvertPCMToWav converts a PCM file to WAV format using ffmpeg.
// Parameters:
// - pcmPath: 경로 PCM 파일 (예: "rx.pcm").
// - wavPath: 변환된 WAV 파일의 경로 (예: "rx.wav").
// - sampleRate: 샘플 레이트 (예: 24000).
// - channels: 채널 수 (예: 1).
func ConvertPCMToWav(pcmPath, wavPath string, sampleRate int, channels int) error {
	// ffmpeg 명령어 구성
	cmd := exec.Command(
		"ffmpeg",
		"-f", "s16le", // PCM 포맷 (16-bit little endian)
		"-ar", fmt.Sprintf("%d", sampleRate), // 샘플 레이트
		"-ac", fmt.Sprintf("%d", channels), // 채널 수
		"-i", pcmPath, // 입력 파일
		wavPath, // 출력 파일
	)

	// 명령어의 표준 출력과 표준 오류를 캡처
	cmd.Stdout = nil // 필요 시 변경 (예: os.Stdout)
	cmd.Stderr = nil // 필요 시 변경 (예: os.Stderr)

	// 명령어 실행
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert %s to %s: %w", pcmPath, wavPath, err)
	}

	return nil
}

// MixPCMToWav mixes two PCM files and saves the result as a WAV file.
func MixPCMToWav(pcmPath1, pcmPath2, wavPath string, sampleRate, channels int) error {
	// ffmpeg 명령어 구성
	cmd := exec.Command(
		"ffmpeg",
		"-f", "s16le", // 첫 번째 PCM 파일의 포맷
		"-ar", fmt.Sprintf("%d", sampleRate), // 샘플 레이트
		"-ac", fmt.Sprintf("%d", channels), // 채널 수
		"-i", pcmPath1, // 첫 번째 PCM 입력 파일
		"-f", "s16le", // 두 번째 PCM 파일의 포맷
		"-ar", fmt.Sprintf("%d", sampleRate), // 샘플 레이트
		"-ac", fmt.Sprintf("%d", channels), // 채널 수
		"-i", pcmPath2, // 두 번째 PCM 입력 파일
		"-filter_complex", "amerge=inputs=2", // 믹스 필터
		"-ac", fmt.Sprintf("%d", channels), // 출력 채널 수
		wavPath, // 출력 WAV 파일 경로
	)

	// 명령어의 표준 출력과 표준 오류를 캡처
	cmd.Stdout = nil // 필요 시 변경 (예: os.Stdout)
	cmd.Stderr = nil // 필요 시 변경 (예: os.Stderr)

	// 명령어 실행
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to mix %s and %s to %s: %w", pcmPath1, pcmPath2, wavPath, err)
	}

	return nil
}
