package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strings"

)

// DefaultProfileTheme menyimpan kombinasi warna untuk profil & background halaman
type DefaultProfileTheme struct {
	Initial        string `json:"initial"`
	PrimaryColor   string `json:"primary_color"`   // warna utama (profil/avatar)
	BackgroundTop  string `json:"background_top"`  // warna gradasi atas
	BackgroundDown string `json:"background_down"` // warna gradasi bawah
	TextColor      string `json:"text_color"`      // otomatis hitam/putih tergantung brightness
}

// GenerateDefaultProfileTheme menghasilkan warna utama, gradasi, dan warna teks
func GenerateDefaultProfileTheme(email string) DefaultProfileTheme {
	namePart := strings.Split(email, "@")[0]
	parts := strings.Split(namePart, ".")

	// --- ambil inisial
	var initial string
	if len(parts) > 1 && len(parts[0]) > 0 && len(parts[1]) > 0 {
		initial = strings.ToUpper(string(parts[0][0]) + string(parts[1][0]))
	} else if len(namePart) >= 2 {
		initial = strings.ToUpper(string(namePart[0]) + string(namePart[1]))
	} else if len(namePart) == 1 {
		initial = strings.ToUpper(string(namePart[0]))
	} else {
		initial = "?"
	}

	// --- warna utama
	hue, sat, light := generateHSLFromEmail(email)
	r, g, b := hslToRgb(hue, sat, light)
	primaryColor := fmt.Sprintf("#%02X%02X%02X", r, g, b)

	// --- gradasi lembut
	rTop, gTop, bTop := hslToRgb(hue, sat*0.9, clamp(light*1.3, 0, 1))
	rDown, gDown, bDown := hslToRgb(hue, sat*0.7, clamp(light*0.5, 0, 1))

	// --- tentukan warna teks berdasarkan kecerahan warna utama
	textColor := getContrastColor(r, g, b)

	return DefaultProfileTheme{
		Initial:        initial,
		PrimaryColor:   primaryColor,
		BackgroundTop:  fmt.Sprintf("#%02X%02X%02X", rTop, gTop, bTop),
		BackgroundDown: fmt.Sprintf("#%02X%02X%02X", rDown, gDown, bDown),
		TextColor:      textColor,
	}
}

// generateHSLFromEmail mengubah hash email → hue/saturation/lightness
func generateHSLFromEmail(email string) (h, s, l float64) {
	hash := md5.Sum([]byte(strings.ToLower(email)))
	hexStr := hex.EncodeToString(hash[:])
	intVal, _ := hex.DecodeString(hexStr[:6])

	sum := int(intVal[0]) + int(intVal[1]) + int(intVal[2])
	h = float64(sum%360) / 360.0
	s = 0.65
	l = 0.55
	return
}

// konversi HSL → RGB
func hslToRgb(h, s, l float64) (r, g, b uint8) {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := l - c/2

	var r1, g1, b1 float64
	switch {
	case 0 <= h && h < 1.0/6.0:
		r1, g1, b1 = c, x, 0
	case 1.0/6.0 <= h && h < 2.0/6.0:
		r1, g1, b1 = x, c, 0
	case 2.0/6.0 <= h && h < 3.0/6.0:
		r1, g1, b1 = 0, c, x
	case 3.0/6.0 <= h && h < 4.0/6.0:
		r1, g1, b1 = 0, x, c
	case 4.0/6.0 <= h && h < 5.0/6.0:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}

	r = uint8((r1 + m) * 255)
	g = uint8((g1 + m) * 255)
	b = uint8((b1 + m) * 255)
	return
}

// hitung kontras warna (putih / hitam)
func getContrastColor(r, g, b uint8) string {
	// rumus luminance berdasarkan sRGB
	luminance := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255
	if luminance > 0.6 {
		return "#000000" // background terang → teks hitam
	}
	return "#FFFFFF" // background gelap → teks putih
}

// clamp menjaga nilai di antara 0 dan 1
func clamp(v, min, max float64) float64 {
	return math.Max(min, math.Min(max, v))
}
