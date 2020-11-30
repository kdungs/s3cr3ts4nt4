package s3cr3ts4nt4

import (
	"encoding/json"
	"fmt"
	"io"
)

type Participant struct {
	Name    string
	Address string
	Pubkey  PublicKey
}

func NewParticipant(name, address string, pubkey PublicKey) *Participant {
	return &Participant{
		Name:    name,
		Address: address,
		Pubkey:  pubkey,
	}
}

func GenerateParticipant(name, address string) (*Participant, *SecretKey, error) {
	keys, err := GenerateKeypair()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate keypair for participant: %w", err)
	}

	return NewParticipant(name, address, keys.Public), &keys.Secret, nil
}

type participantInfo struct {
	Name    string
	Address string
}

type participantPayload struct {
	EncryptedInfo []byte
	Pubkey        PublicKey
}

func (p Participant) WriteEncryped(w io.Writer, public PublicKey) error {
	infoBytes, err := json.Marshal(participantInfo{
		Name:    p.Name,
		Address: p.Address,
	})
	if err != nil {
		return fmt.Errorf("unable to serialize participant info: %w", err)
	}

	encryptedInfo, err := Encrypt(public, infoBytes)
	if err != nil {
		return fmt.Errorf("unable to encrypt participant info: %w", err)
	}

	payload := participantPayload{
		EncryptedInfo: encryptedInfo,
		Pubkey:        p.Pubkey,
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("unable to serialize participant payload: %w", err)
	}
	return nil
}

func ReadEncryptedParticipant(r io.Reader, secret SecretKey) (*Participant, error) {
	var payload participantPayload
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, fmt.Errorf("unable to deserialize participant payload: %w", err)
	}

	infoBytes, err := Decrypt(secret, payload.EncryptedInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt participant info: %w", err)
	}

	var info participantInfo
	if err := json.Unmarshal(infoBytes, &info); err != nil {
		return nil, fmt.Errorf("unable to deserialize participant info: %w", err)
	}

	return &Participant{
		Name:    info.Name,
		Address: info.Address,
		Pubkey:  payload.Pubkey,
	}, nil
}
