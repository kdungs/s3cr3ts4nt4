package cmd

import (
	"errors"
	"fmt"
	"os"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

func addHost(rootCmd *cobra.Command, cli *s3cr3ts4nt4.CLI) {
	hostCmd := &cobra.Command{
		Use:   "host",
		Short: "host a gift exchange",
	}
	var identity string
	hostCmd.PersistentFlags().StringVarP(
		&identity,
		"identity",
		"i",
		"host",
		"identity to use for the host",
	)

	hostNewCmd := &cobra.Command{
		Use:   "new",
		Short: "start a new gift exchange",
		Long: `Start a new gift exchange.

If no host.id file exists, this will create a new host identity. It will then
export the host public key so it can be sent to participants.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.HostNew(identity)
		},
	}
	hostCmd.AddCommand(hostNewCmd)

	var outdir string
	hostRunCmd := &cobra.Command{
		Use:   "run [participant payload files]",
		Short: "run a gift exchange",
		Args: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(outdir); !os.IsNotExist(err) {
				return fmt.Errorf("directory %s already exists", outdir)
			}
			if len(args) < 2 {
				return errors.New("at least two participants are required")
			}
			for _, fname := range args {
				if _, err := os.Stat(fname); os.IsNotExist(err) {
					return fmt.Errorf("%s does not exist", fname)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.HostRun(identity, outdir, args)
		},
	}
	hostRunCmd.Flags().StringVarP(
		&outdir,
		"outdir",
		"o",
		"results",
		"directory that results will be written to; must not exist and will be created",
	)
	hostCmd.AddCommand(hostRunCmd)

	rootCmd.AddCommand(hostCmd)
}
