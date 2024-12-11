package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	LogLevel       = logrus.InfoLevel
	RmsThresholdDb = -50.0 // -50 dBFS 이하일 경우 무음으로 간주
	UseZCR         = false
	ZcrThreshold   = 0.15 // ZCR이 15% 이하일 경우 무음으로 간주
	SystemPrompt   = func() string {
		file, err := os.ReadFile("config/tutor_prompt.txt")
		if err != nil {
			return ""
		}

		return string(file)
	}
)
