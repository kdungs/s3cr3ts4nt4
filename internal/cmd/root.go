package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "s3cr3ts4nt4",
	Short: "🎅 cryptographic gift exchange",
	Long: `s3cr3ts4nt4
🎅 cryptographic gift exchange 🎅
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
