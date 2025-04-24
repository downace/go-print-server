# Print Server

A simple HTTP printing server.

Main purpose is to allow fast and silent printing from any device in local network

Built using [Wails 2](https://wails.io/) and Vue with [Quasar](https://quasar.dev/).

## Supported Platforms

- Linux - supported using `lp` and `lpstat`
- Windows - supported using `wmic` and embedded [`SumatraPDF`](https://www.sumatrapdfreader.org/)
- macOS - not supported (PR's are welcome)

## Server API

> Currently, API is very limited, just allowing to list printers and print PDF files without any options.
> Feel free to file an issue if you need more features or options.

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

## Roadmap

- [x] Build releases with GitHub Actions
- [x] Windows support
- [x] Printing from URL
- [x] HTTPS support
- [x] CORS (via custom response headers)
- [ ] ...
