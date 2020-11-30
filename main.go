package main

import (
	"io"
	"log"
	"sync"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
)

func run() error {
	keys, err := s3cr3ts4nt4.GenerateKeypair()
	if err != nil {
		return err
	}

	p1, _, err := s3cr3ts4nt4.GenerateParticipant("Hans", "Hans sein Haus")
	if err != nil {
		return err
	}
	p2, _, err := s3cr3ts4nt4.GenerateParticipant("Dieter", "Dieter sein Haus")
	if err != nil {
		return err
	}
	p3, _, err := s3cr3ts4nt4.GenerateParticipant("Klaus", "Klaus sein Haus")
	if err != nil {
		return err
	}

	santa := s3cr3ts4nt4.NewSanta(*keys)

	addParticipant := func(p s3cr3ts4nt4.Participant) error {
		r, w := io.Pipe()
		var wg sync.WaitGroup
		var goerr error
		wg.Add(1)
		go func() {
			if goerr = santa.AddEncryptedParticipant(r); goerr != nil {
				return
			}
			wg.Done()
		}()
		if err := p.WriteEncryped(w, keys.Public); err != nil {
			return err
		}
		w.Close()
		wg.Wait()
		return goerr
	}

	addParticipants := func(ps ...s3cr3ts4nt4.Participant) error {
		for _, p := range ps {
			if err := addParticipant(p); err != nil {
				return err
			}
		}
		return nil
	}

	if err := addParticipants(*p1, *p2, *p3); err != nil {
		return err
	}

	res, err := santa.Run()
	if err != nil {
		return err
	}
	log.Printf("%v", res)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}
