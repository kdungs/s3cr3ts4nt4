package s3cr3ts4nt4

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Santa struct {
	Participants []Participant
	keys         KeyPair
}

func NewSanta(keypair KeyPair) *Santa {
	return &Santa{
		Participants: make([]Participant, 0),
		keys:         keypair,
	}
}

func SantaFromSecret(r io.Reader) (*Santa, error) {
	var sec SecretKey
	if err := json.NewDecoder(r).Decode(&sec); err != nil {
		return nil, fmt.Errorf("unable to deserialize secret key: %w", err)
	}
	return NewSanta(KeyPairFromSecretKey(sec)), nil
}

func (s *Santa) AddParticipant(p Participant) error {
	for _, op := range s.Participants {
		if op.Name == p.Name {
			return fmt.Errorf("participant %s already exists", p.Name)
		}
	}
	s.Participants = append(s.Participants, p)
	return nil
}

func (s Santa) Run() (map[string][]byte, error) {
	l := len(s.Participants)
	if l < 2 {
		return nil, fmt.Errorf("cannot do a gift exchange with %d users", l)
	}
	indices, err := DerangedIndices(l)
	if err != nil {
		return nil, fmt.Errorf("unable to derange indices: %w", err)
	}
	result := make(map[string][]byte, l)
	for senderIdx, recipientIdx := range indices {
		sender := s.Participants[senderIdx]
		recipient := s.Participants[recipientIdx]

		var buf bytes.Buffer
		if err := recipient.WriteEncryped(&buf, sender.Pubkey); err != nil {
			return nil, fmt.Errorf("unable to encrypt participant %s: %w", recipient.Name, err)
		}

		result[sender.Name] = buf.Bytes()
	}

	return result, nil
}

func (s Santa) StoreKeys(w io.Writer) error {
	if err := json.NewEncoder(w).Encode(s.keys); err != nil {
		return fmt.Errorf("unable to serialize Santa's keys: %w", err)
	}
	return nil
}

func (s *Santa) AddEncryptedParticipant(r io.Reader) error {
	p, err := ReadEncryptedParticipant(r, s.keys.Secret)
	if err != nil {
		return fmt.Errorf("unable to add encrypted participant: %w", err)
	}
	if err := s.AddParticipant(*p); err != nil {
		return fmt.Errorf("unable to add participant: %w", err)
	}
	return nil
}
