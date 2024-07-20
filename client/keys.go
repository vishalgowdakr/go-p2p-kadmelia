package client

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

const keyFile = "node_key.b64"

func LoadOrCreatePrivateKey() (crypto.PrivKey, error) {
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		// Create new key
		priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
		if err != nil {
			return nil, err
		}
		// Save to file
		if err := savePrivateKey(priv); err != nil {
			return nil, err
		}
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

func CreateIDFromPrivateKey(privateKey crypto.PrivKey) (string, error) {
	// 1. Get the public key from the private key
	publicKey := privateKey.GetPublic()

	// 2. Convert the public key to bytes
	publicKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	// 3. Hash the public key using SHA-256
	hash := sha256.Sum256(publicKeyBytes)

	// 4. Convert the hash to a hexadecimal string
	id := hex.EncodeToString(hash[:])

	return id, nil
}

/*
// Use when creating your libp2p host
priv, err := LoadOrCreatePrivateKey()
if err != nil {
    // handle error
}
host, err := libp2p.New(
    libp2p.Identity(priv),
    // other options...
)
*/
