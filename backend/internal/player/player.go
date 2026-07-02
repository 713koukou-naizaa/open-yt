package player

import (
	"fmt"
	"os"
	"os/exec"

	"open-yt/internal/config"
)

func Play(videoURL string, cfg config.PlayerConfiguration) error {
	MPVArgs := []string{}

	if cfg.YTDLFormat != "" {
		MPVArgs = append(MPVArgs, fmt.Sprintf("--ytdl-format=%s", cfg.YTDLFormat))
	}
	if cfg.Volume != nil {
		MPVArgs = append(MPVArgs, fmt.Sprintf("--volume=%d", *cfg.Volume))
	}
	if cfg.Fullscreen != "" {
		MPVArgs = append(MPVArgs, fmt.Sprintf("--fullscreen=%s", cfg.Fullscreen))
	}
	if cfg.WindowMaximized != "" {
		MPVArgs = append(MPVArgs, fmt.Sprintf("--window-maximized=%s", cfg.WindowMaximized))
	}
	if cfg.KeepOpen != "" {
		MPVArgs = append(MPVArgs, fmt.Sprintf("--keep-open=%s", cfg.KeepOpen))
	}
	if cfg.ForceWindow != "" {
		MPVArgs = append(MPVArgs, fmt.Sprintf("--force-window=%s", cfg.ForceWindow))
	}

	cmd := exec.Command(cfg.Command, append(MPVArgs, videoURL)...)

	// Pipe command's stdout and stderr to current process's stdout and stderr
	// so user can see output from player
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run will block until player is closed
	cmdRunError := cmd.Run()
	// If error is ExitError, player was closed by user
	// Not fatal application error, so return nil
	_, cmdRunIsExitError := cmdRunError.(*exec.ExitError)
	if cmdRunIsExitError {
		return nil
	}
	// For other errors (e.g., command not found), should report them
	return cmdRunError
}
