package audiomanager

import (
	"context"
	"github.com/gordonklaus/portaudio"
	"github.com/sirupsen/logrus"
	"openai-realtime/pkg/config"
	"os"
)

var log = func() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel)
	return log
}()

type Manager struct {
	DeviceController *Controller
	SampleRate       int
	VolumeThresh     float32

	outputChan chan []int16 // audio data channel to speaker
	errorChan  chan error
}

func NewManager(inputDevice *portaudio.DeviceInfo, outputDevice *portaudio.DeviceInfo, volumeThreshold float32) (*Manager, error) {

	return &Manager{
		DeviceController: NewController(inputDevice, outputDevice, volumeThreshold),
		SampleRate:       int(inputDevice.DefaultSampleRate),
		VolumeThresh:     volumeThreshold,
		errorChan:        make(chan error),
	}, nil
}

func (m *Manager) Start(ctx context.Context) error {
	log.Info("Audio manager started")
	return m.DeviceController.On(ctx)
}

func (m *Manager) Close() {
	m.DeviceController.Off()
	log.Info("Audio manager stopped")
}
