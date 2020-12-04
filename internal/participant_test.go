package s3cr3ts4nt4

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestParticipantCanWriteAndReadEncrypted(t *testing.T) {
	is := is.New(t)

	p, _, err := NewParticipant(
		"Bilbo Baggins",
		"A hole in the ground\nThe Shire",
	)
	is.NoErr(err)

	identity, err := GenerateIdentity()
	is.NoErr(err)

	var buf bytes.Buffer
	is.NoErr(p.WriteEncryped(&buf, identity.Public()))

	decrypted, err := ReadEncryptedParticipant(&buf, *identity)
	is.NoErr(err)
	is.Equal(p, decrypted)
}
