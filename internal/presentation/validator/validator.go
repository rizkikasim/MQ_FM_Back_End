package validator

import (
	"errors"
	"fmt"
	"strings"

)

// ValidateRequiredFields memastikan semua field yang wajib diisi tidak kosong.
func ValidateRequiredFields(fields map[string]string) error {
	var missing []string
	for key, value := range fields {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return errors.New(fmt.Sprintf("field wajib diisi: %s", strings.Join(missing, ", ")))
	}
	return nil
}

// ValidatePasswordLength memastikan password memenuhi panjang minimal.
func ValidatePasswordLength(password string, min int) error {
	if len(strings.TrimSpace(password)) < min {
		return fmt.Errorf("password minimal %d karakter", min)
	}
	return nil
}
