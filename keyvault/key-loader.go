package keyvault

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"log"
	"os"
)

func LoadEncryptionKey() ([]byte, error) {
	log.Println("🔍 [KeyVault] Starting LoadEncryptionKey()")

	vaultURL := os.Getenv("KEYVAULT_URL")
	secretName := os.Getenv("ENCRYPTION_KEY_NAME")

	if vaultURL == "" || secretName == "" {
		log.Println("❌ [KeyVault] KEYVAULT_URL or ENCRYPTION_KEY_NAME not set")
		return nil, fmt.Errorf("missing KEYVAULT_URL or ENCRYPTION_KEY_NAME")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Printf("❌ [KeyVault] Failed to create credential: %v", err)
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}
	log.Println("✅ [KeyVault] Credential created")

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		log.Printf("❌ [KeyVault] Failed to create client: %v", err)
		return nil, fmt.Errorf("failed to create keyvault client: %w", err)
	}
	log.Println("✅ [KeyVault] Client created")

	resp, err := client.GetSecret(context.Background(), secretName, "", nil)
	if err != nil {
		log.Printf("❌ [KeyVault] Failed to get secret: %v", err)
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	log.Println("✅ [KeyVault] Secret retrieved")

	keyBytes, err := base64.StdEncoding.DecodeString(*resp.Value)
	if err != nil {
		log.Printf("❌ [KeyVault] Failed to decode base64 key: %v", err)
		return nil, fmt.Errorf("failed to decode base64 key: %w", err)
	}
	log.Println("✅ [KeyVault] Key decoded successfully")

	return keyBytes, nil
}
