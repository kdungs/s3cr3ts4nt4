package cmd

import (
	"errors"
	"fmt"
	"os"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

func addDecrypt(rootCmd *cobra.Command) {
	var myIdentityFile string
	decryptCmd := &cobra.Command{
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
			return s3cr3ts4nt4.NewCli().Decrypt(myIdentityFile, args[0])
		},
	}
	decryptCmd.Flags().StringVarP(
		&myIdentityFile,
		"identity",
		"i",
		"me",
		"your identity file",
	)
	rootCmd.AddCommand(decryptCmd)
}
