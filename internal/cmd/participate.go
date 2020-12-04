package cmd

import (
	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

func addParticipate(rootCmd *cobra.Command, cli *s3cr3ts4nt4.CLI) {
	var identity string
	var hostkeyFile string
	var name string
	var address string

	participateCmd := &cobra.Command{
		Use:   "participate",
		Short: "participate in a gift exchange",
		Long: `Participate in a gift exchange.

This will create a file called "YOUR NAME.in" which only the host of the
exchange can decrypt. Send it to them and wait for them to send you an
encrypted file containing your recipient.  Please make sure you don't lose your
identity file (e.g. me.id).
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.Participate(
				hostkeyFile,
				identity,
				name,
				address,
			)
		},
	}
	participateCmd.Flags().StringVarP(
		&identity,
		"identity",
		"i",
		"me",
		"your identity; if it doesn't exist, it will be created",
	)
	participateCmd.Flags().StringVarP(
		&hostkeyFile,
		"hostkey",
		"k",
		"host.pub",
		"host public key",
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
		"your address; you can use a single-quoted string with newlines",
	)
	participateCmd.MarkFlagRequired("address")

	rootCmd.AddCommand(participateCmd)
}
