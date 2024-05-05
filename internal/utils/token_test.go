package utils

import (
	"log"
	"testing"
	"time"

	"github.com/enzo-gbd/GBA/configs"
	"github.com/stretchr/testify/require"
)

func TestGenerateTokenAndValidateToken(t *testing.T) {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load environment variables: %v", err)
	}

	payload := "testUser"
	ttl := time.Hour

	tokenString, err := GenerateToken(ttl, payload, config.AccessTokenPrivateKey)
	require.NoError(t, err, "La génération du token ne devrait pas échouer")

	retrievedPayload, err := ValidateToken(tokenString, config.AccessTokenPublicKey)
	require.NoError(t, err, "La validation du token ne devrait pas échouer")
	require.Equal(t, payload, retrievedPayload, "Le payload récupéré devrait correspondre au payload d'origine")
}
