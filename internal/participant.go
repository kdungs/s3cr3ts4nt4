package s3cr3ts4nt4

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (p Participant) WriteEncryped(w io.Writer, pubkey NaclKey) error {
	data, err := json.Marshal(p)
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
