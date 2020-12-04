package cmd

import (
	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

func addParticipate(rootCmd *cobra.Command) {
	var identity string
	var hostkeyFile string
	var name string
	var address string
	var outfile string

	participateCmd := &cobra.Command{
		Use:   "participate",
		Short: "participate in a gift exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3cr3ts4nt4.NewCli().Participate(
				identity,
				hostkeyFile,
				name,
				address,
				outfile,
			)
		},
	}
	participateCmd.Flags().StringVarP(
		&identity,
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
