package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "golang.org/x/image/webp"
)

// scanPhotos walks the given directory and categorizes image files by orientation.
// It reads image headers to determine dimensions without fully decoding pixels.
func scanPhotos(dir string) (landscape, portrait []string) {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !validExts[ext] {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		config, _, err := image.DecodeConfig(file)
		if err != nil {
			return nil
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil
		}

		if config.Width >= config.Height {
			landscape = append(landscape, absPath)
		} else {
			portrait = append(portrait, absPath)
		}
		return nil
	})
	return
}

func main() {
	application := app.New()
	window := application.NewWindow("Digital Picture Frame")
	window.SetFullScreen(true)

	landscape, portrait := scanPhotos("./photos")

	if len(landscape) == 0 && len(portrait) == 0 {
		window.SetContent(widget.NewLabel("No photos found in ./photos"))
		window.ShowAndRun()
		return
	}

	// Pre-allocate image widgets and containers so we can reuse them across slides.
	landscapeImg := canvas.NewImageFromFile("")
	landscapeImg.FillMode = canvas.ImageFillContain
	landscapeContainer := container.NewStack(landscapeImg)

	portraitImgLeft := canvas.NewImageFromFile("")
	portraitImgLeft.FillMode = canvas.ImageFillContain
	portraitImgRight := canvas.NewImageFromFile("")
	portraitImgRight.FillMode = canvas.ImageFillContain
	portraitContainer := container.NewGridWithColumns(2, portraitImgLeft, portraitImgRight)

	// showNext picks a random orientation and displays the corresponding photo(s).
	// Portrait mode requires at least 2 photos to show side-by-side; otherwise landscape is used.
	showNext := func() {
		canLandscape := len(landscape) > 0
		canPortrait := len(portrait) >= 2

		usePortrait := false
		if canLandscape && canPortrait {
			usePortrait = rand.IntN(2) == 0
		} else if canPortrait {
			usePortrait = true
		}

		if usePortrait {
			// Pick two distinct random portrait photos.
			firstIdx := rand.IntN(len(portrait))
			secondIdx := rand.IntN(len(portrait) - 1)
			if secondIdx >= firstIdx {
				secondIdx++
			}
			fyne.Do(func() {
				portraitImgLeft.File = portrait[firstIdx]
				portraitImgRight.File = portrait[secondIdx]
				portraitImgLeft.Refresh()
				portraitImgRight.Refresh()
				window.SetContent(portraitContainer)
			})
		} else {
			photoIdx := rand.IntN(len(landscape))
			fyne.Do(func() {
				landscapeImg.File = landscape[photoIdx]
				landscapeImg.Refresh()
				window.SetContent(landscapeContainer)
			})
		}
	}

	// Start the slideshow: display immediately, then swap every 10 seconds.
	go func() {
		showNext()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			showNext()
		}
	}()

	window.ShowAndRun()
}
