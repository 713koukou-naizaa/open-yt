package youtube

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Fetches the playlists saved by the logged-in user.
// Requires yt-dlp to be configured with cookies for a logged-in YouTube account.
func GetPlaylists(paginationThreshold int, YTDLPCommand string, browser string) ([]YTPlaylist, error) {
	YTDLPArgs := []string{
		"--flat-playlist",
		"--dump-json",
		"--playlist-end",
		fmt.Sprintf("%d", paginationThreshold),
	}
	if browser != "" {
		YTDLPArgs = append(YTDLPArgs, "--cookies-from-browser", browser)
	}
	YTDLPArgs = append(YTDLPArgs, feedPlaylistsURL)

	var YTPlaylists []YTPlaylist
	processor := func(line []byte) error {
		var currentYTDLPPlaylist YTDLPPlaylist
		if err := json.Unmarshal(line, &currentYTDLPPlaylist); err != nil {
			return err
		}
		YTPlaylists = append(YTPlaylists, newYTPlaylistFromYTDLPPlaylist(currentYTDLPPlaylist))
		return nil
	}

	if err := ytdlpExecutor(YTDLPCommand, YTDLPArgs, processor); err != nil {
		return nil, fmt.Errorf("playlists fetch failed: %w", err)
	}

	return YTPlaylists, nil
}

// Fetches videos from a selected playlist.
func GetPlaylistVideos(playlistURL string, paginationThreshold int, YTDLPCommand string, browser string) ([]YTVideo, error) {
	YTDLPArgs := []string{
		"--flat-playlist",
		"--dump-json",
		"--playlist-end",
		fmt.Sprintf("%d", paginationThreshold),
	}
	if browser != "" {
		YTDLPArgs = append(YTDLPArgs, "--cookies-from-browser", browser)
	}
	YTDLPArgs = append(YTDLPArgs, playlistURL)

	var YTVideos []YTVideo
	processor := func(line []byte) error {
		var currentYTDLPVideo YTDLPVideo
		if err := json.Unmarshal(line, &currentYTDLPVideo); err != nil {
			return err
		}
		YTVideos = append(YTVideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
		return nil
	}

	if err := ytdlpExecutor(YTDLPCommand, YTDLPArgs, processor); err != nil {
		return nil, fmt.Errorf("playlist videos fetch failed: %w", err)
	}

	return YTVideos, nil
}

func newYTPlaylistFromYTDLPPlaylist(YTDLPPlaylist YTDLPPlaylist) YTPlaylist {
	playlistURL := YTDLPPlaylist.WebpageURL
	if playlistURL == "" {
		playlistURL = YTDLPPlaylist.URL
	}
	if playlistURL != "" && strings.HasPrefix(playlistURL, "/") {
		playlistURL = BaseURL + playlistURL
	}
	if playlistURL == "" && YTDLPPlaylist.ID != "" {
		playlistURL = BaseURL + "/playlist?list=" + YTDLPPlaylist.ID
	}

	return YTPlaylist{
		ID:         YTDLPPlaylist.ID,
		Title:      YTDLPPlaylist.Title,
		URL:        playlistURL,
		VideoCount: YTDLPPlaylist.PlaylistCount,
	}
}
