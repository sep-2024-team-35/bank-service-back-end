package config

import (
	"github.com/sep-2024-team-35/bank-servce-back-end/keyvault"
	"log"
)

var EncryptionKey []byte

func LoadEncryptionKeyGlobally() {
	var err error
	EncryptionKey, err = keyvault.LoadEncryptionKey()
	if err != nil {
		log.Fatalf("❌ Failed to load encryption key: %v", err)
	}
	log.Println("🔐 Encryption key loaded successfully")
}
