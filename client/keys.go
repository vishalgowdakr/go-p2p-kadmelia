package client

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

const keyFile = "node_key.b64"

func loadOrCreatePrivateKey() (crypto.PrivKey, error) {
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		// Create new key
		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
		if err != nil {
			return nil, err
		}

		// Save to file
		savePrivateKey(priv)
		return priv, nil
	}

	// Load existing key
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	keyData, err := base64.StdEncoding.DecodeString(string(keyBytes))
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(keyData)
}

func savePrivateKey(priv crypto.PrivKey) error {
	keyBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return err
	}
	keyB64 := base64.StdEncoding.EncodeToString(keyBytes)
	return os.WriteFile(keyFile, []byte(keyB64), 0600)
}

/* // Use when creating your libp2p host
priv, err := loadOrCreatePrivateKey()
if err != nil {
    // handle error
}
host, err := libp2p.New(
    libp2p.Identity(priv),
    // other options...
) */
