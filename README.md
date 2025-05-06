# Print Server

A simple HTTP printing server.

Main purpose is to allow fast and silent printing from any device in local network

Two variants available: CLI app and desktop app with GUI.

GUI app is built using [Wails 2](https://wails.io/) and Vue 3 with [Quasar](https://quasar.dev/).

## Supported Platforms

- Linux - supported using `lp` and `lpstat`
- Windows - supported using `wmic` and embedded [`SumatraPDF`](https://www.sumatrapdfreader.org/)
- macOS - not supported (PR's are welcome)

## Usage

Download suitable binary from [Releases](https://github.com/downace/go-print-server/releases) and start it.

For CLI version, use `*-cli` binary, e.g. `print-server-linux-v0.7.0-cli`.  
Use `-help` flag to see available options.

You can also [build manually from sources](#development)

## Server API

- `GET /printers` - get list of available printers
   ```shell
   curl http://127.0.0.1:8888/printers
   ```
   ```json
   {"printers": [{"name":"Brother_MFC_L2700DN_series"},{"name":"PDF"}]}
   ```
- `POST /print-pdf` - print PDF file
   ```shell
   curl --header 'Content-Type: application/pdf' --data-binary /path/to/file.pdf http://127.0.0.1:8888/print-pdf?printer=Brother_MFC_L2700DN_series
   ```
- `POST /print-pdf-url` - print PDF file from URL
   ```shell
   curl http://127.0.0.1:8888/print-pdf-url?printer=Brother_MFC_L2700DN_series&url=https%3A%2F%2Fpdfobject.com%2Fpdf%2Fsample.pdf
   ```
- `POST /print-url` - print any file from URL
   > Loaded file is converted to PDF with [Rod](https://go-rod.github.io) which uses Chromium by default.
   > First call to this method may take some time as Chromium needs to be loaded

   Query params: see `PrintFromUrlQuery` in `internal/server/handlers.go`

   ```shell
   curl http://127.0.0.1:8888/print-pdf-url?printer=Brother_MFC_L2700DN_series&url=https%3A%2F%2Fhttpstat.us%2F&pages=2-7
   ```

## Development

### GUI app

First, [Install Wails CLI](https://wails.io/docs/gettingstarted/installation#installing-wails)

Run Development mode:

```shell
wails dev
```

Build the app:

```shell
wails build
```

### CLI

Use `cli` tag to compile and run CLI app:

```shell
go build -tags cli
```

```shell
go run -tags cli .
```
