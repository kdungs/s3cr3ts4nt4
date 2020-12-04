package s3cr3ts4nt4

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/box"
)

const (
	naclkeyBytes = 32
)

type NaclKey [naclkeyBytes]byte

func ImportNaclKey(r io.Reader) (*NaclKey, error) {
	var key NaclKey
	n, err := r.Read(key[:])
	if err != nil {
		return nil, err
	}
	if n != naclkeyBytes {
		return nil, fmt.Errorf(
			"expected to read %d bytes, got %d",
			naclkeyBytes,
			n,
		)
	}
	return &key, nil
}

func (key NaclKey) Export(w io.Writer) error {
	n, err := w.Write(key[:])
	if err != nil {
		return err
	}
	if n != naclkeyBytes {
		return fmt.Errorf(
			"expected to write %d bytes, wrote %d",
			naclkeyBytes,
			n,
		)
	}
	return nil
}

func (key NaclKey) Encrypt(msg []byte) ([]byte, error) {
	var out []byte
	pubkey := [naclkeyBytes]byte(key)
	enc, err := box.SealAnonymous(out, msg, &pubkey, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt message: %w", err)
	}
	return enc, nil
}

func NaclDecrypt(msg []byte, pub NaclKey, priv NaclKey) ([]byte, error) {
	var out []byte
	pubkey := [naclkeyBytes]byte(pub)
	privkey := [naclkeyBytes]byte(priv)
	dec, ok := box.OpenAnonymous(out, msg, &pubkey, &privkey)
	if !ok {
		return nil, errors.New("unable to decrypt message")
	}
	return dec, nil
}

type Identity struct {
	public  NaclKey
	private NaclKey
}

func GenerateIdentity() (*Identity, error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Identity{
		public:  *pub,
		private: *priv,
	}, nil
}

func ImportIdentity(r io.Reader) (*Identity, error) {
	pub, err := ImportNaclKey(r)
	if err != nil {
		return nil, fmt.Errorf("unable to import public key: %w", err)
	}
	priv, err := ImportNaclKey(r)
	if err != nil {
		return nil, fmt.Errorf("unable to import private key: %w", err)
	}
	return &Identity{
		public:  *pub,
		private: *priv,
	}, nil
}

func (i *Identity) Export(w io.Writer) error {
	if err := i.public.Export(w); err != nil {
		return fmt.Errorf("unable to export public key: %w", err)
	}
	if err := i.private.Export(w); err != nil {
		return fmt.Errorf("unable to export private key: %w", err)
	}
	return nil
}

func (i *Identity) Public() NaclKey {
	return i.public
}

func (i *Identity) Decrypt(msg []byte) ([]byte, error) {
	dec, err := NaclDecrypt(msg, i.public, i.private)
	if err != nil {
		return nil, err
	}
	return dec, nil
}
