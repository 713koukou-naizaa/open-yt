<br />
<div align="center">
  <!-- <a href="https://github.com/713koukou-naizaa/open-yt">
    <img src="assets/open-yt-logo.png" alt="Open-yt logo" width="80" height="80">
  </a> -->

<h1 align="center">Open-yt</h1>

  <p align="center">
    A modern, lightweight YouTube client without ads/tracking.
    <br />
    <a href="https://github.com/713koukou-naizaa/open-yt/issues">Report Bug</a>
    ·
    <a href="https://github.com/713koukou-naizaa/open-yt/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#configuration">Configuration</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

## About The Project

![Open-yt Demo](./assets/open-yt-demo.gif)

Open-yt is a modern, lightweight YouTube client removing ads/tracking from your YouTube experience. It's designed to be easy to use, only the main YouTube features, distraction-free.

It uses `yt-dlp` to fetch video information and `mpv` to stream the content, providing a fast and efficient viewing experience.

### Built With

*   Go
*   yt-dlp
*   mpv
*   Bubble Tea
*   Viper

## Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

You need to have the following software installed:
*   **Go** (version 1.21 or later)
*   **yt-dlp**
*   **mpv** (or another video player, configurable)

For features like accessing your home feed and subscriptions, you'll need to have `yt-dlp` configured with your YouTube cookies. The easiest way is to have a compatible browser with your cookies available.

### Installation

1.  Clone the repo
    ```sh
    git clone https://github.com/713koukou-naizaa/open-yt.git
    ```
2.  Build the application
    ```sh
    cd open-yt/backend
    make
    ```
3.  (Optional) Move the binary to a directory in your `PATH` for easy access.
    ```sh
    sudo mv open-yt /usr/local/bin/
    ```

## Usage

Run `open-yt` without any arguments to launch the interactive TUI menu:
```sh
open-yt
```

You can also use the raw search and play features directly from the command line:
*   Search for videos:
    ```sh
    open-yt search "your search query"
    ```
*   Play a video from a URL or ID:
    ```sh
    open-yt play "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
    ```

## Configuration

Open-yt can be configured via a YAML file located at `~/.config/open-yt/config.yaml`.

The application will create a default configuration if one is not found. You can customize the `yt-dlp` command, the video player and its arguments, and more.

See the `config.go` file for all available options.

## Roadmap

- [x] Search videos
- [x] Play video from link/ID
- [x] Browse subscriptions feed
- [x] Browse individual channel uploads (videos/streams)
- [x] Browse home feed
- [x] Interactive TUI with dynamic filtering
- [x] External configuration file
- [ ] Pagination in search results
- [ ] Playlist support
- [ ] Display thumbnails in the terminal
- [ ] GUI Frontend (Qt/C++?)

See the open issues for a full list of proposed features (and known issues).

## Contributing

Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

## Acknowledgments

*   yt-dlp for the heavy lifting of fetching YouTube content.
*   mpv for being a superb, scriptable video player.
*   The Charm team for their fantastic Go libraries that make building beautiful TUIs a joy.
