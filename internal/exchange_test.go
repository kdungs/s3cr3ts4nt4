package s3cr3ts4nt4

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestExchange(t *testing.T) {
	is := is.New(t)

	id, err := GenerateIdentity()
	is.NoErr(err)
	exchange := NewExchange(*id)

	infos := map[string]string{
		"James Jameson":    "123 Some Street\nABC 123 Some town\nEngland",
		"Hans Hansen":      "Einestraße 23\n12345 Einestadt\nDeutschland",
		"Giacomo Gianluca": "Via Esempio 1\nLorem Citta\nItalia",
		"Testy McTestface": "Tester road\nLoch Ness\nScottland",
	}
	identities := make(map[string]Identity, len(infos))
	for name, addr := range infos {
		p, pid, err := NewParticipant(name, addr)
		is.NoErr(err)
		is.NoErr(exchange.AddParticipant(*p))
		identities[name] = *pid
	}

	// Make sure we cannot add a participant with an existing name.
	badp, _, err := NewParticipant("Hans Hansen", "Some address")
	is.NoErr(err)
	err = exchange.AddParticipant(*badp)
	is.True(err != nil)

	// Ensure exchange runs.
	res, err := exchange.Run()
	is.NoErr(err)
	is.Equal(4, len(res))

	// Ensure that each participant can decrypt _exactly one_ result, and that
	// everybody appears as a recipient _exactly once_.
	recipientCount := make(map[string]int, len(infos))
	for gifterName, payload := range res {
		gifter, ok := identities[gifterName]
		is.True(ok)
		buf := bytes.NewBuffer(payload)
		recipient, err := ReadEncryptedParticipant(buf, gifter)
		is.NoErr(err)
		is.True(recipient.Name != gifterName)
		is.Equal(recipient.Address, infos[recipient.Name])
		recipientCount[recipient.Name]++
	}
	for name := range infos {
		count, ok := recipientCount[name]
		is.True(ok)
		is.Equal(1, count)
	}

}
