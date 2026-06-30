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

- Search videos from search term
- Play video from link/ID

### Second MVP - Advanced CLI version

- Pagination in searches
- Get videos from subscriptions feed
- Get subscriptions
- Search videos from search term from subscription

### Third MVP - Full CLI version

abc

### Fourth MVP - abc
