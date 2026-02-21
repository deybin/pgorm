package configs

import (
	"os"
	"strings"
)

// GetConfig busca una variable de entorno.
// Si no existe, devuelve un valor por defecto.
func get(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

// IsValidationEnabled verifica si el usuario activó las validaciones en su .env
func IsValidationEnabled() bool {
	// Ejemplo: PGORM_VALIDATION_STRICT=true
	return strings.ToLower(get("PGORM_VALIDATION_STRICT", "false")) == "true"
}

func KeyCrypto() []byte {
	v := os.Getenv("ENV_KEY_CRYPTO")
	if v == "" {
		return nil
	}
	return []byte(v)
}
