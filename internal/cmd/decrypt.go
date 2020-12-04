package cmd

import (
	"errors"
	"fmt"
	"os"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

func addDecrypt(rootCmd *cobra.Command, cli *s3cr3ts4nt4.CLI) {
	var identity string
	decryptCmd := &cobra.Command{
		Use:   "decrypt recipientFile",
		Short: "decrypt recipient information",
		Long: `Decrypt your recipient's information.

Once your host has sent you your recipient file ("YOUR NAME.out"), you can
decrypt it using this command. The program will then show you the name and
address of your recipient whom you can then send a gift.
`,
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
			return cli.Decrypt(identity, args[0])
		},
	}
	decryptCmd.Flags().StringVarP(
		&identity,
		"identity",
		"i",
		"me",
		"your identity",
	)
	rootCmd.AddCommand(decryptCmd)
}
