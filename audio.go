package main

import (
	"fmt"
	"os/exec"
)

func ProcessAudio(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg not found %v", err)
	}

	cmd = exec.Command(
		"ffmpeg",
		"-y",
		"-i",
		inputPath,
		"-filter:a",
		"atempo=2",
		outputPath,
	)

	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to speed up audio using ffmpeg %v", err)
	}

	return nil
}
