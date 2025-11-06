package utils

import (
	authadminentity "mqfm_backend/internal/domain/entities/auth/admin"
	"mqfm_backend/internal/helper"

)

// GenerateDefaultProfileThemeWrapper
// biar service tetap bisa pakai tipe ProfileTheme dari entity
func GenerateDefaultProfileTheme(identifier string) authadminentity.ProfileTheme {
	h := helper.GenerateDefaultProfileTheme(identifier)
	return authadminentity.ProfileTheme{
		Initial:        h.Initial,
		PrimaryColor:   h.PrimaryColor,
		BackgroundTop:  h.BackgroundTop,
		BackgroundDown: h.BackgroundDown,
		TextColor:      h.TextColor,
	}
}
