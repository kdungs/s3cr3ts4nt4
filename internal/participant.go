package s3cr3ts4nt4

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
)

type Participant struct {
	Name    string
	Address string
	Pubkey  NaclKey
}

func NewParticipant(name string, address string) (*Participant, *Identity, error) {
	identity, err := GenerateIdentity()
	if err != nil {
		return nil, nil, err
	}
	return &Participant{
		Name:    name,
		Address: address,
		Pubkey:  identity.Public(),
	}, identity, nil
}

func randomPadding() ([]byte, error) {
	// TODO(kdungs): Make functions that use rand.Reader take an io.Reader
	// instead.
	rnd := rand.Reader
	const minPadding = 128
	const maxPadding = 1024

	pSizeBig, err := rand.Int(rnd, big.NewInt(maxPadding-minPadding))
	if err != nil {
		return nil, err
	}
	pSize := pSizeBig.Int64() + minPadding

	bs := make([]byte, pSize)
	if _, err := rnd.Read(bs); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return bs, nil
}

func (p Participant) WriteEncryped(w io.Writer, pubkey NaclKey) error {
	padding, err := randomPadding()
	if err != nil {
		return fmt.Errorf("unable to generate random padding: %w", err)
	}
	payload := struct {
		Participant
		Padding []byte
	}{
		Participant: p,
		Padding:     padding,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to serialize participant info: %w", err)
	}
	enc, err := pubkey.Encrypt(data)
	if err != nil {
		return fmt.Errorf("unable to encrypt participant info: %w", err)
	}
	if _, err := w.Write(enc); err != nil {
		return fmt.Errorf("unable to write encrypted participant info: %w", err)
	}
	return nil
}

func ReadEncryptedParticipant(r io.Reader, identity Identity) (*Participant, error) {
	msg, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read: %w", err)
	}
	data, err := identity.Decrypt(msg)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt participant info: %w", err)
	}
	var p Participant
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("unable to deserialize participant info: %w", err)
	}
	return &p, nil
}
