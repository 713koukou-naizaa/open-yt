package player

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"open-yt/internal/config"
)

func Play(videoURL string, cfg config.PlayerConfiguration) error {
	return play(videoURL, cfg, nil)
}

func PlayPlaylist(videoURLs []string, startIndex int, cfg config.PlayerConfiguration, cookiesFromBrowser string) error {
	if startIndex < 0 || startIndex >= len(videoURLs) {
		return fmt.Errorf("playlist start index %d is out of range", startIndex)
	}

	playlistFile, err := os.CreateTemp("", "open-yt-*.m3u")
	if err != nil {
		return fmt.Errorf("failed to create temporary playlist: %w", err)
	}
	playlistFileName := playlistFile.Name()
	defer os.Remove(playlistFileName)

	var playlistContent strings.Builder
	playlistContent.WriteString("#EXTM3U\n")
	for _, videoURL := range videoURLs[startIndex:] {
		if strings.ContainsAny(videoURL, "\r\n") {
			playlistFile.Close()
			return fmt.Errorf("playlist video URL contains a newline")
		}
		playlistContent.WriteString(videoURL)
		playlistContent.WriteString("\n")
	}

	if _, err := playlistFile.WriteString(playlistContent.String()); err != nil {
		playlistFile.Close()
		return fmt.Errorf("failed to write temporary playlist: %w", err)
	}
	if err := playlistFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary playlist: %w", err)
	}

	additionalArgs := []string{}
	if cookiesFromBrowser != "" {
		additionalArgs = append(additionalArgs, fmt.Sprintf("--ytdl-raw-options-append=cookies-from-browser=%s", cookiesFromBrowser))
	}
	return play(playlistFileName, cfg, additionalArgs)
}

func play(videoURL string, cfg config.PlayerConfiguration, additionalArgs []string) error {
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
	MPVArgs = append(MPVArgs, additionalArgs...)

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
