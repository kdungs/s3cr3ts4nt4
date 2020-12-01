package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

var myIdentityFile string

func init() {
	decryptCmd.Flags().StringVarP(
		&myIdentityFile,
		"identity",
		"i",
		"me",
		"your identity file",
	)
	rootCmd.AddCommand(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt recipientFile",
	Short: "decrypt a recipients address that was encrypted for you",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("recipient file has to be provided")
		}
		recipientFile := args[0]
		if _, err := os.Stat(recipientFile); os.IsNotExist(err) {
			return fmt.Errorf(
				"recipient file %s does not exist: %w",
				recipientFile,
				err,
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎅 Decrypting recipient payload 🎅")
		// Load identity
		fh, err := os.Open(myIdentityFile)
		if err != nil {
			return fmt.Errorf("unable to open identity %s: %w", myIdentityFile, err)
		}
		defer fh.Close()
		var sec s3cr3ts4nt4.SecretKey
		if err := json.NewDecoder(fh).Decode(&sec); err != nil {
			return fmt.Errorf("unable to deserialize identity: %w", err)
		}

		// Decrypt
		recipientFile := args[0]
		f, err := os.Open(recipientFile)
		if err != nil {
			return fmt.Errorf("unable to open %s: %w", recipientFile, err)
		}
		defer f.Close()
		p, err := s3cr3ts4nt4.ReadEncryptedParticipant(f, sec)
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
🎁 🎅 🎁 🎅 🎁 🎅 🎁 🎅
`,
			p.Name,
			p.Address,
		)

		return nil
	},
}
