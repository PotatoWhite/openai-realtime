package audiomanager

import (
	"bufio"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"os"
	"strconv"
	"strings"
)

func selectDevice(deviceType string) *portaudio.DeviceInfo {
	devices, err := portaudio.Devices()
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	fmt.Println("Available devices:")
	for i, device := range devices {
		var maxChannels int
		if deviceType == "input" {
			maxChannels = device.MaxInputChannels
		} else {
			maxChannels = device.MaxOutputChannels
		}

		if maxChannels > 0 {
			fmt.Printf("%d: %s (MaxInputChannels: %d, MaxOutputChannels: %d, DefaultSampleRate: %.0f)\n",
				i, device.Name, device.MaxInputChannels, device.MaxOutputChannels, device.DefaultSampleRate)
		}
	}

	fmt.Printf("Enter the number of the %s device to use: ", deviceType)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	input = strings.TrimSpace(input)
	deviceIndex, err := strconv.Atoi(input)
	if err != nil || deviceIndex < 0 || deviceIndex >= len(devices) {
		log.Fatalf("Invalid device number")
	}

	return devices[deviceIndex]
}

func SelectInputDevice() *portaudio.DeviceInfo {
	return selectDevice("input")
}

func SelectOutputDevice() *portaudio.DeviceInfo {
	return selectDevice("output")
}
