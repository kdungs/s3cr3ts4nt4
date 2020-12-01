package s3cr3ts4nt4

// Just a few convenience wrappers around crypto.

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"hash"
)

const (
	numBits = 4096
)

func defaultHash() hash.Hash {
	return sha256.New()
}

type PublicKey struct {
	Key rsa.PublicKey
}

type SecretKey struct {
	Key rsa.PrivateKey
}

type KeyPair struct {
	Public PublicKey
	Secret SecretKey
}

func GenerateKeypair() (*KeyPair, error) {
	key, err := rsa.GenerateKey(rand.Reader, numBits)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Public: PublicKey{key.PublicKey},
		Secret: SecretKey{*key},
	}, nil
}

func Encrypt(public PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(defaultHash(), rand.Reader, &public.Key, msg, nil)
}

func Decrypt(secret SecretKey, payload []byte) ([]byte, error) {
	return rsa.DecryptOAEP(defaultHash(), rand.Reader, &secret.Key, payload, nil)
}

func KeyPairFromSecretKey(sec SecretKey) KeyPair {
	return KeyPair{
		Public: PublicKey{sec.Key.PublicKey},
		Secret: sec,
	}
}
