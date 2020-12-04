package s3cr3ts4nt4

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Cli struct{}

func NewCli() *Cli {
	return &Cli{}
}

func (c *Cli) HostNew(identity string) error {
	fmt.Println("🎅 Starting a new gift exchange 🎅")
	keypair, err := GenerateKeypair()
	if err != nil {
		return fmt.Errorf("unable to generate host keypair: %w", err)
	}
	parts := map[string]interface{}{
		fmt.Sprintf("%s.pub", identity): keypair.Public,
		fmt.Sprintf("%s.sec", identity): keypair.Secret,
	}
	// First, check if any of the files already exist. In that case, we
	// don't want to overwrite them as we'd risk losing the key for a live
	// gift exchange.
	for fname, _ := range parts {
		if _, err := os.Stat(fname); !os.IsNotExist(err) {
			return fmt.Errorf("file %s already exists", fname)
		}
	}
	// Serialize the keys to distinct files.
	for fname, part := range parts {
		fh, err := os.Create(fname)
		if err != nil {
			return fmt.Errorf(
				"unable to open file %s for writing: %w",
				fname,
				err,
			)
		}
		defer fh.Close()
		if err := json.NewEncoder(fh).Encode(part); err != nil {
			return fmt.Errorf(
				"unable to serialize host key to %s: %w",
				fname,
				err,
			)
		}
	}
	fmt.Printf("Done.\nSend %s.pub to your participants.\n", identity)
	return nil
}

func (c *Cli) HostRun(identity, resultDir string, participants []string) error {
	fmt.Println("🎅 Running a gift exchange 🎅")
	fmt.Printf("Participants: \n - %s\n", strings.Join(participants, "\n - "))

	// Create secret santa instance.
	secretKeyFile := fmt.Sprintf("%s.sec", identity)
	fh, err := os.Open(secretKeyFile)
	if err != nil {
		return fmt.Errorf(
			"unable to open secret file %s: %w",
			secretKeyFile,
			err,
		)
	}
	defer fh.Close()
	santa, err := SantaFromSecret(fh)
	if err != nil {
		return fmt.Errorf(
			"unable to create santa from secret file %s: %w",
			secretKeyFile,
			err,
		)
	}

	// Add participants
	for _, participantFile := range participants {
		fh, err := os.Open(participantFile)
		if err != nil {
			return fmt.Errorf(
				"unable to open participant payload file %s: %w",
				participantFile,
				err,
			)
		}
		defer fh.Close()
		if err := santa.AddEncryptedParticipant(fh); err != nil {
			fmt.Errorf(
				"unable to add participant from file %s: %w",
				participantFile,
				err,
			)
		}
	}

	// Create output directory
	if err := os.MkdirAll(resultDir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf(
			"unable to create directory %s: %w",
			resultDir,
			err,
		)
	}

	// Run the spiel.
	mapping, err := santa.Run()
	if err != nil {
		fmt.Errorf("unable to assign gifts: %w", err)
	}
	for name, payload := range mapping {
		fname := fmt.Sprintf("%s/%s", resultDir, name)
		fh, err := os.Create(fname)
		if err != nil {
			return fmt.Errorf("unable to create %s: %w", fname, err)
		}
		defer fh.Close()
		if _, err := fh.Write(payload); err != nil {
			return fmt.Errorf(
				"unable to write payload to %s: %w",
				fname,
				err,
			)
		}
	}

	fmt.Printf(`
Done.
Distribute the files in %s to your participants.

🎁 Enjoy your gifts! 🎁
`,
		resultDir,
	)
	return nil
}

func (c *Cli) Decrypt(identityFile string, recipientFile string) error {
	fmt.Println("🎅 Decrypting recipient payload 🎅")
	// Load identity
	fident, err := os.Open(identityFile)
	if err != nil {
		return fmt.Errorf("unable to open identity %s: %w", identityFile, err)
	}
	defer fident.Close()
	var sec SecretKey
	if err := json.NewDecoder(fident).Decode(&sec); err != nil {
		return fmt.Errorf("unable to deserialize identity: %w", err)
	}

	// Decrypt
	frecipient, err := os.Open(recipientFile)
	if err != nil {
		return fmt.Errorf("unable to open %s: %w", recipientFile, err)
	}
	defer frecipient.Close()
	p, err := ReadEncryptedParticipant(frecipient, sec)
	if err != nil {
		return fmt.Errorf("unable to decrypt recipient: %w", err)
	}

	fmt.Printf(`
Done.
Please send a gift to

%s
%s

🎁 🎅 🎁 🎅 🎁 🎅 🎁 🎅
 Happy gift exchange!
🎅 🎁 🎅 🎁 🎅 🎁 🎅 🎁
`,
		p.Name,
		p.Address,
	)

	return nil
}

func (c *Cli) Participate(
	identityFile string,
	hostkeyFile string,
	name string,
	address string,
	outfile string,
) error {
	fmt.Println("🎅 Generating participant payload 🎅")
	// Load identity if it exists; otherwise generate and write to file
	var identity KeyPair
	if _, err := os.Stat(identityFile); !os.IsNotExist(err) {
		fh, err := os.Open(identityFile)
		if err != nil {
			return fmt.Errorf("unable to open %s: %w", identityFile, err)
		}
		defer fh.Close()
		var sec SecretKey
		if err := json.NewDecoder(fh).Decode(&sec); err != nil {
			return fmt.Errorf("unable to deserialize identity: %w", err)
		}
		identity = KeyPairFromSecretKey(sec)
	} else {
		kp, err := GenerateKeypair()
		if err != nil {
			return fmt.Errorf("unable to generate keypair: %w", err)
		}
		fh, err := os.Create(identityFile)
		if err != nil {
			return fmt.Errorf("unable to create %s: %w", identityFile, err)
		}
		defer fh.Close()
		if err := json.NewEncoder(fh).Encode(kp.Secret); err != nil {
			return fmt.Errorf("unable to serialize identity: %w", err)
		}
		identity = *kp
	}

	// Load host public key
	fh, err := os.Open(hostkeyFile)
	if err != nil {
		return fmt.Errorf(
			"unable to open host key file %s: %w",
			hostkeyFile,
			err,
		)
	}
	defer fh.Close()
	var hostkey PublicKey
	if err := json.NewDecoder(fh).Decode(&hostkey); err != nil {
		return fmt.Errorf("unable to deserialize host key: %w", err)
	}

	// Create and encrypt payload.
	p := NewParticipant(name, address, identity.Public)
	if outfile == "" {
		outfile = fmt.Sprintf("%s.out", name)
	}
	fh, err = os.Create(outfile)
	if err != nil {
		return fmt.Errorf("unable to create %s: %w", outfile, err)
	}
	defer fh.Close()
	if err := p.WriteEncryped(fh, hostkey); err != nil {
		return fmt.Errorf("unable to encrypt payload: %w", err)
	}

	fmt.Printf(`
Done.
Give %s to your host.
They will send you another user's encrypted information, later.
`,
		outfile,
	)
	return nil
}
