package s3cr3ts4nt4

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestIdentityCanExportAndImport(t *testing.T) {
	is := is.New(t)

	id, err := GenerateIdentity()
	is.NoErr(err)

	var buf bytes.Buffer
	is.NoErr(id.Export(&buf))

	imported, err := ImportIdentity(&buf)
	is.NoErr(err)
	is.Equal(id, imported)
}

func TestIdentityCanEncryptAndDecrypt(t *testing.T) {
	is := is.New(t)

	id, err := GenerateIdentity()
	is.NoErr(err)

	msg := []byte("Hello, Alice!")
	enc, err := id.Public().Encrypt(msg)
	is.NoErr(err)

	dec, err := id.Decrypt(enc)
	is.NoErr(err)
	is.Equal(msg, dec)
}
