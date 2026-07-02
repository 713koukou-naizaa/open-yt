# Open-yt

Open-yt is a modern lightweight youtube client available both in CLI and GUI mode.

## Components

The Open-yt client is composed of a backend in Go, and 2 frontends : a CLI one in Go, and a GUI one in Qt with C++ 

## Inner-workings

- Backend uses yt-dlp to get the content only (videos, subscriptions...), and MPV to stream videos
- Content exposed via API
- CLI Frontend allows features like search, watch etc... but from CLI
- Both Frontends get content from the same API

## Plan

Implement little-by-little each components from simplest to more complex ones

### First MVP - Minimal CLI version

#### Features

- [x] Search videos from search term
- [x] Play video from link/ID

### Second MVP - Advanced CLI version

- [x] Get videos from subscriptions feed
- [x] Get subscriptions
- [x] Get videos from subscriptions
- [x] Get lives from subscriptions
- [x] Search videos from search term from subscription
- [x] Dynamic filter in interactive list
- [x] Configuration via external file

### Third MVP - Full CLI version

- [ ] Pagination in searches

### Fourth MVP - abc
