package helper

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	colorextractor "github.com/marekm4/color-extractor"

)

// ImageProfileTheme menyimpan warna dari foto profil
type ImageProfileTheme struct {
	PrimaryColor   string `json:"primary_color"`
	BackgroundTop  string `json:"background_top"`
	BackgroundDown string `json:"background_down"`
	TextColor      string `json:"text_color"`
}

// GenerateImageProfileTheme membuat tema warna dari file image (khusus foto profil)
func GenerateImageProfileTheme(imagePath string) (ImageProfileTheme, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return ImageProfileTheme{}, fmt.Errorf("gagal membuka file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return ImageProfileTheme{}, fmt.Errorf("gagal decode gambar: %v", err)
	}

	colors := colorextractor.ExtractColors(img)
	if len(colors) == 0 {
		return ImageProfileTheme{}, fmt.Errorf("tidak bisa ekstrak warna dari gambar")
	}

	// warna dominan (langsung ambil color.Color)
	dominant := colors[0]
	r, g, b, _ := dominant.RGBA()
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

	primary := fmt.Sprintf("#%02X%02X%02X", r8, g8, b8)
	text := getContrastColorFromImage(r8, g8, b8)

	// gradasi sederhana
	rTop := clampColor(r8, 1.15)
	gTop := clampColor(g8, 1.15)
	bTop := clampColor(b8, 1.15)

	rDown := clampColor(r8, 0.85)
	gDown := clampColor(g8, 0.85)
	bDown := clampColor(b8, 0.85)

	return ImageProfileTheme{
		PrimaryColor:   primary,
		BackgroundTop:  fmt.Sprintf("#%02X%02X%02X", rTop, gTop, bTop),
		BackgroundDown: fmt.Sprintf("#%02X%02X%02X", rDown, gDown, bDown),
		TextColor:      text,
	}, nil
}

func clampColor(value uint8, factor float64) uint8 {
	v := float64(value) * factor
	if v > 255 {
		v = 255
	}
	if v < 0 {
		v = 0
	}
	return uint8(v)
}

func getContrastColorFromImage(r, g, b uint8) string {
	luminance := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255
	if luminance > 0.6 {
		return "#000000"
	}
	return "#FFFFFF"
}
