package pgorm

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetKey_PrivateCrypto() []byte {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error configuración de variables de entorno")
	}
	key := os.Getenv("ENV_KEY_CRYPTO")
	return []byte(key)
}
