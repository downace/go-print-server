# Icons

Converting SVG icons (using Inkscape) to PNG:

```shell
inkscape build/appicon.svg -w 256 -h 256 -o build/appicon.png
```

Converting PNG icon to ICO for Windows (using ImageMagick)

```shell
convert build/appicon.png build/windows/icon.ico
```
