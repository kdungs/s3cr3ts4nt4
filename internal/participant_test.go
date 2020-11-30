package s3cr3ts4nt4

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestParticipantEncryption(t *testing.T) {
	is := is.New(t)

	keys, err := GenerateKeypair()
	is.NoErr(err)

	p, _, err := GenerateParticipant(
		"Bilbo Baggins",
		"A hole in the ground\nThe Shire",
	)
	is.NoErr(err)

	var buf bytes.Buffer
	is.NoErr(p.WriteEncryped(&buf, keys.Public))

	decrypted, err := ReadEncryptedParticipant(&buf, keys.Secret)
	is.NoErr(err)
	is.Equal(p, decrypted)
}
