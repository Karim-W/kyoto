# Kyoto

TUI application to display the currently playing song on your Spotify account.

> Named after Phoebe Bridgers' Kyoto because its a great song

## Pre-Requisites

- Spotify Account
- Spotify Application
- MacOS
- Golang

## Setup Process

### 1. Install Kyoto

```bash
go install github.com/karim-w/kyoto@latest
```

### 2. Create a Spotify Application

The spotify application needed so that kyoto can request the currently playing
song from the Spotify API.

You can create a Spotify application [here](https://developer.spotify.com/dashboard/applications).

For reference on how to create a Spotify application, you can follow the
[Spotify Documentation](https://developer.spotify.com/documentation/web-api/concepts/apps).

#### 2.1 Set the Redirect URI

The redirect URI is the URI that Spotify will redirect to after the user has
authorized the application.

The Redirect URI should be set to `http://localhost:9499/callback`.

> its very important that the redirect URI is set to `http://localhost:9499/callback`
> because that is the url that kyoto listens to for the authorization code.

## Usage

### 1. Authorize Kyoto

```bash
kyoto login
```

This will open a browser window and prompt you to login to your Spotify account
and authorize the application.

After a successful login kyoto will presist the authorization token in
`~/.kyoto` file

### 2. Start Kyoto

```bash
kyoto
```

This will start kyoto and display the currently playing song on your Spotify
account.
![kyoto](./docs/kyoto.png)

## Roadmap

- [ ] Add support for Linux
- [ ] Dynamic Backgrounds
- [ ] Lyric Viewer

## License

BSD 3-Clause License

## Contributing

Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.

## Acknowledgements

- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)

## Author

Karim-W
