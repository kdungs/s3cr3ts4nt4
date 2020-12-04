package cmd

import (
	"fmt"
	"os"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "s3cr3ts4nt4",
		Short: "🎅 cryptographic gift exchange",
		Long: `🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁
     _____          _____ _       _  _         _   _  _
 ___|___ /  ___ _ _|___ /| |_ ___| || |  _ __ | |_| || |
/ __| |_ \ / __| '__||_ \| __/ __| || |_| '_ \| __| || |_
\__ \___) | (__| |  ___) | |_\__ \__   _| | | | |_|__   _|
|___/____/ \___|_| |____/ \__|___/  |_| |_| |_|\__|  |_|

🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅🎁🎅

Cryptographic gift exchange
github.com/kdungs/s3cr3ts4nt4
`,
	}
	cli := s3cr3ts4nt4.NewCLI()
	addHost(rootCmd, cli)
	addParticipate(rootCmd, cli)
	addDecrypt(rootCmd, cli)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
