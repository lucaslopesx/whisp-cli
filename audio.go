package main

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
)

type ProcessingOptions struct {
	InputPath  string
	OutputPath string
	Velocity   float64
}

func ProcessAudio(opts ProcessingOptions) error {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg not found %v", err)
	}

	atempoFilter := buildAtempoFilter(opts.Velocity)
	fmt.Print(atempoFilter)

	cmd = exec.Command(
		"ffmpeg",
		"-y",
		"-i",
		opts.InputPath,
		"-filter:a",
		atempoFilter,
		opts.OutputPath,
	)

	err = cmd.Run()

	if err != nil {
		fmt.Print(opts.InputPath)
		fmt.Print(opts.OutputPath)
		return fmt.Errorf("failed to speed up audio using ffmpeg %v", err)
	}

	return nil
}

func buildAtempoFilter(velocity float64) string {
	var parts []string
	for velocity > 2.0 {
		parts = append(parts, "2.0")
		velocity /= 2.0
	}

	velocity = math.Round(velocity*1000) / 1000
	parts = append(parts, fmt.Sprintf("%.3f", velocity))

	return "atempo=" + strings.Join(parts, ",atempo=")
}
