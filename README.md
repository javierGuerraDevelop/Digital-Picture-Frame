# Digital Picture Frame

A fullscreen slideshow app built with Go and [Fyne](https://fyne.io/).

It scans a `./photos` directory for images (jpg, png, webp), categorizes them by orientation, and displays them in a rotating slideshow every 10 seconds. Landscape photos are shown full-screen; portrait photos are shown two side-by-side.

## Usage

Place your images in a `photos/` directory next to the binary, then run:

```sh
go run .
```

## Build

```sh
go build -o digital-picture-frame .
```

## Requirements

- Go 1.23+
- OpenGL support (Fyne dependency)
