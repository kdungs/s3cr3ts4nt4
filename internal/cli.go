package s3cr3ts4nt4

import (
	"fmt"
	"os"
	"path"
)

type CLI struct{}

func NewCLI() *CLI {
	return &CLI{}
}

func loadIdentity(identity string) (*Identity, error) {
	identityFile := fmt.Sprintf("%s.id", identity)
	fh, err := os.Open(identityFile)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to open %s for reading: %w",
			identityFile,
			err,
		)
	}
	defer fh.Close()
	id, err := ImportIdentity(fh)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func loadOrGenerateIdentity(identity string) (*Identity, error) {
	identityFile := fmt.Sprintf("%s.id", identity)
	if _, err := os.Stat(identityFile); err == nil {
		// file exists: read and return Identity
		return loadIdentity(identity)
	} else if os.IsNotExist(err) {
		// file does not exist: generate identity and write to file
		id, err := GenerateIdentity()
		if err != nil {
			return nil, fmt.Errorf("unable to generate identity: %w", err)
		}
		fh, err := os.Create(identityFile)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to open %s for writing: %w",
				identityFile,
				err,
			)
		}
		defer fh.Close()
		if err := id.Export(fh); err != nil {
			return nil, err
		}
		return id, nil
	} else {
		return nil, err
	}
}

func (c *CLI) HostNew(identity string) error {
	fmt.Println("🎅 Starting a new gift exchange 🎅")
	id, err := loadOrGenerateIdentity(identity)
	if err != nil {
		return err
	}

	publicFile := fmt.Sprintf("%s.pub", identity)
	fh, err := os.Create(publicFile)
	if err != nil {
		return fmt.Errorf("unable to open %s for writing: %w", publicFile, err)
	}
	defer fh.Close()
	if err := id.Public().Export(fh); err != nil {
		return fmt.Errorf("unable to export host public key: %w", err)
	}

	fmt.Printf(`Done.
Give %s to your participants.
`,
		publicFile,
	)
	return nil
}

func (c *CLI) HostRun(
	identity string,
	outdir string,
	participantsFiles []string,
) error {
	fmt.Println("🎅 Running a gift exchange 🎅")
	id, err := loadIdentity(identity)
	if err != nil {
		return err
	}
	exchange := NewExchange(*id)

	for _, pf := range participantsFiles {
		fh, err := os.Open(pf)
		if err != nil {
			return fmt.Errorf(
				"unable to open participant file %s: %w",
				pf,
				err,
			)
		}
		defer fh.Close()
		p, err := ReadEncryptedParticipant(fh, *id)
		if err != nil {
			return fmt.Errorf(
				"unable to decrypt participant file %s: %w",
				pf,
				err,
			)
		}
		if err := exchange.AddParticipant(*p); err != nil {
			return fmt.Errorf("unable to add participant: %w", err)
		}
	}

	if err := os.MkdirAll(outdir, os.ModePerm|os.ModeDir); err != nil {
		return fmt.Errorf("unable to create %s: %w", outdir, err)
	}

	res, err := exchange.Run()
	if err != nil {
		return fmt.Errorf("unable to run gift exchange: %w", err)
	}
	for name, payload := range res {
		fname := path.Join(outdir, fmt.Sprintf("%s.out", name))
		fh, err := os.Create(fname)
		if err != nil {
			return fmt.Errorf("unable to open %s for writing: %w", fname, err)
		}
		defer fh.Close()
		if _, err := fh.Write(payload); err != nil {
			return fmt.Errorf("unable to write payload to %s: %w", fname, err)
		}
	}

	fmt.Printf(`Done.
Give each participant their respective file from %s.
`,
		outdir,
	)
	return nil
}

func (c *CLI) Participate(
	hostkey string,
	identity string,
	name string,
	address string,
) error {
	// Read host public key
	fh, err := os.Open(hostkey)
	if err != nil {
		return fmt.Errorf("unable to open host key file %s: %w", hostkey, err)
	}
	defer fh.Close()
	hostpub, err := ImportNaclKey(fh)
	if err != nil {
		return fmt.Errorf("unable to read host public key: %w", err)
	}

	id, err := loadOrGenerateIdentity(identity)
	if err != nil {
		return err
	}
	// TODO: refactor NewParticipant
	p := &Participant{
		Name:    name,
		Address: address,
		Pubkey:  id.Public(),
	}

	fname := fmt.Sprintf("%s.in", p.Name)
	outf, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("unable to open %s for writing: %w", fname, err)
	}
	defer outf.Close()
	if err := p.WriteEncryped(outf, *hostpub); err != nil {
		return fmt.Errorf(
			"unable to write participant file %s: %w",
			fname,
			err,
		)
	}
	fmt.Printf(`Done.
Send %s to your host.
`,
		fname,
	)
	return nil
}

func (c *CLI) Decrypt(
	identity string,
	recipientFile string,
) error {
	id, err := loadIdentity(identity)
	if err != nil {
		return err
	}

	fh, err := os.Open(recipientFile)
	if err != nil {
		return fmt.Errorf(
			"unable to open recipient file %s: %w",
			recipientFile,
			err,
		)
	}
	defer fh.Close()
	r, err := ReadEncryptedParticipant(fh, *id)
	if err != nil {
		return fmt.Errorf(
			"unable to read recipient file %s: %w",
			recipientFile,
			err,
		)
	}

	fmt.Printf(`Done.
Please send a gift to

%s
%s

🎁 🎅 🎁 🎅 🎁 🎅 🎁 🎅
 Happy gift exchange!
🎅 🎁 🎅 🎁 🎅 🎁 🎅 🎁
`,
		r.Name,
		r.Address,
	)
	return nil
}
