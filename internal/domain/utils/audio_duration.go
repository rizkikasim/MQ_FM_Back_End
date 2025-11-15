package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

)

type ffprobeOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func GetAudioDuration(filePath string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		filePath,
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("ffprobe error: %v", err)
	}

	var result ffprobeOutput
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return 0, err
	}

	if result.Format.Duration == "" {
		return 0, fmt.Errorf("duration tidak ditemukan")
	}

	var seconds float64
	fmt.Sscanf(result.Format.Duration, "%f", &seconds)

	return int(seconds), nil
}
