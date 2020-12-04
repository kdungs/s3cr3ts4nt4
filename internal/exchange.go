package s3cr3ts4nt4

import (
	"bytes"
	"fmt"
)

type Exchange struct {
	identity     Identity
	participants []Participant
}

func NewExchange(id Identity) *Exchange {
	return &Exchange{
		identity:     id,
		participants: make([]Participant, 0),
	}
}

func (s *Exchange) AddParticipant(p Participant) error {
	for _, op := range s.participants {
		if op.Name == p.Name {
			return fmt.Errorf("participant %s already exists", p.Name)
		}
	}
	s.participants = append(s.participants, p)
	return nil
}

func (s *Exchange) Run() (map[string][]byte, error) {
	n := len(s.participants)
	if n < 2 {
		return nil, fmt.Errorf(
			"cannot do a gift exchange with %d participants",
			n,
		)
	}
	indices, err := DerangedIndices(n)
	if err != nil {
		return nil, fmt.Errorf("unable to derange indices: %w", err)
	}
	result := make(map[string][]byte, n)
	for senderIdx, recipientIdx := range indices {
		sender := s.participants[senderIdx]
		recipient := s.participants[recipientIdx]
		var buf bytes.Buffer
		if err := recipient.WriteEncryped(&buf, sender.Pubkey); err != nil {
			return nil, fmt.Errorf(
				"unable to encrypt participant %s: %w",
				recipient.Name,
				err,
			)
		}
		result[sender.Name] = buf.Bytes()
	}

	return result, nil
}
