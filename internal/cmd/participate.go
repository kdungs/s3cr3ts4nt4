package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

var identityFile string
var hostkeyFile string
var name string
var address string
var outfile string

func init() {
	participateCmd.Flags().StringVarP(
		&identityFile,
		"identity",
		"i",
		"me",
		"your identity file; if it doesn't exist, it will be create on the fly",
	)
	participateCmd.Flags().StringVarP(
		&hostkeyFile,
		"hostkey",
		"k",
		"host.pub",
		"host public key to use for the gift exchange",
	)
	participateCmd.Flags().StringVarP(
		&name,
		"name",
		"n",
		"",
		"your name",
	)
	participateCmd.MarkFlagRequired("name")
	participateCmd.Flags().StringVarP(
		&address,
		"address",
		"a",
		"",
		"your address, you can use a quoted string with \\n",
	)
	participateCmd.MarkFlagRequired("address")
	participateCmd.Flags().StringVarP(
		&outfile,
		"outfile",
		"o",
		"",
		"the file to write your payload to; if not set it will use your name",
	)

	rootCmd.AddCommand(participateCmd)
}

var participateCmd = &cobra.Command{
	Use:   "participate",
	Short: "participate in a gift exchange",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎅 Generating participant payload 🎅")
		// Load identity if it exists; otherwise generate and write to file
		var identity s3cr3ts4nt4.KeyPair
		if _, err := os.Stat(identityFile); !os.IsNotExist(err) {
			fh, err := os.Open(identityFile)
			if err != nil {
				return fmt.Errorf("unable to open %s: %w", identityFile, err)
			}
			defer fh.Close()
			var sec s3cr3ts4nt4.SecretKey
			if err := json.NewDecoder(fh).Decode(&sec); err != nil {
				return fmt.Errorf("unable to deserialize identity: %w", err)
			}
			identity = s3cr3ts4nt4.KeyPairFromSecretKey(sec)
		} else {
			kp, err := s3cr3ts4nt4.GenerateKeypair()
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
		var hostkey s3cr3ts4nt4.PublicKey
		if err := json.NewDecoder(fh).Decode(&hostkey); err != nil {
			return fmt.Errorf("unable to deserialize host key: %w", err)
		}

		// Create and encrypt payload.
		p := s3cr3ts4nt4.NewParticipant(name, address, identity.Public)
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
	},
}
