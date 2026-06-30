package player

import (
	"fmt"
	"os"
	"os/exec"

	"open-yt/internal/config"
)

func Play(video string, cfg config.Player) error {
	args := []string{}

	if cfg.YTDLFormat != "" {
		args = append(args, fmt.Sprintf("--ytdl-format=%s", cfg.YTDLFormat))
	}
	if cfg.Volume != nil {
		args = append(args, fmt.Sprintf("--volume=%d", *cfg.Volume))
	}
	if cfg.Fullscreen != "" {
		args = append(args, fmt.Sprintf("--fullscreen=%s", cfg.Fullscreen))
	}
	if cfg.WindowMaximized != "" {
		args = append(args, fmt.Sprintf("--window-maximized=%s", cfg.WindowMaximized))
	}
	if cfg.KeepOpen != "" {
		args = append(args, fmt.Sprintf("--keep-open=%s", cfg.KeepOpen))
	}
	if cfg.ForceWindow != "" {
		args = append(args, fmt.Sprintf("--force-window=%s", cfg.ForceWindow))
	}

	cmd := exec.Command(cfg.Command, append(args, video)...)

	// Pipe command's stdout and stderr to current process's stdout and stderr
	// so user can see output from player
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run will block until player is closed
	return fmt.Errorf("failed to run player command: %w", cmd.Run())
}
