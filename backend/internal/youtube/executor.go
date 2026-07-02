package youtube

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Handles common logic for executing yt-dlp command
// and processing its line-by-line JSON output
func ytdlpExecutor(ytdlpCommand string, args []string, processLine func(line []byte) error) error {
	cmd := exec.Command(ytdlpCommand, args...)

	stdout, cmdStdoutPipeError := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", cmdStdoutPipeError)
	}

	stderr, cmdStderrPiepError := cmd.StderrPipe()
	if cmdStderrPiepError != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", cmdStderrPiepError)
	}

	cmdStartError := cmd.Start()
	if cmdStartError != nil {
		return fmt.Errorf("failed to start yt-dlp command: %w", cmdStartError)
	}

	// Read stderr in background to capture any error messages,
	// especially if command fails early
	stderrBytes, _ := io.ReadAll(stderr)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Bytes()
		// If line not empty, try to process it
		if len(line) > 0 {
			lineProcessError := processLine(line)
			if lineProcessError != nil {
				// Processor function can decide to log and continue
				// Useful for non-JSON warning lines from yt-dlp
				fmt.Fprintf(os.Stderr, "yt-dlp warning: %s\n", string(line))
			}
		}
	}

	scannerError := scanner.Err()
	if scannerError != nil {
		return fmt.Errorf("error reading yt-dlp output: %w", scannerError)
	}

	cmdWaitError := cmd.Wait()
	if cmdWaitError != nil {
		// Combine command exit error with captured stderr for a comprehensive message
		return fmt.Errorf("yt-dlp command failed: %w\n%s", cmdWaitError, string(stderrBytes))
	}

	return nil
}