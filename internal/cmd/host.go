package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	s3cr3ts4nt4 "github.com/kdungs/s3cr3ts4nt4/internal"
	"github.com/spf13/cobra"
)

// host new
var identity string

// host run
var secretKeyFile string
var resultDir string

func init() {
	hostNewCmd.Flags().StringVarP(
		&identity,
		"identity",
		"i",
		"host",
		"identity to use for the host, will create a .pub and .sec file",
	)
	hostCmd.AddCommand(hostNewCmd)

	hostRunCmd.Flags().StringVarP(
		&secretKeyFile,
		"secret",
		"s",
		"host.sec",
		"host secret key to use for the gift exchange",
	)
	hostRunCmd.Flags().StringVarP(
		&resultDir,
		"outdir",
		"o",
		"results",
		"directory that results will be written to; must not exist and will be created",
	)
	hostCmd.AddCommand(hostRunCmd)

	rootCmd.AddCommand(hostCmd)
}

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "host a gift exchange",
}

var hostNewCmd = &cobra.Command{
	Use:   "new",
	Short: "start a new gift exchange",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎅 Starting a new gift exchange 🎅")
		keypair, err := s3cr3ts4nt4.GenerateKeypair()
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
	},
}

var hostRunCmd = &cobra.Command{
	Use:   "run [participant payload files]",
	Short: "run a gift exchange",
	Args: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(resultDir); !os.IsNotExist(err) {
			return fmt.Errorf("directory %s already exists", resultDir)
		}
		if len(args) < 2 {
			return errors.New("at least two participants are required")
		}
		for _, fname := range args {
			if _, err := os.Stat(fname); os.IsNotExist(err) {
				return fmt.Errorf("payload file %s does not exist", fname)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎅 Running a gift exchange 🎅")
		fmt.Printf("Participants: \n - %s\n", strings.Join(args, "\n - "))

		// Create secret santa instance.
		fh, err := os.Open(secretKeyFile)
		if err != nil {
			return fmt.Errorf(
				"unable to open secret file %s: %w",
				secretKeyFile,
				err,
			)
		}
		defer fh.Close()
		santa, err := s3cr3ts4nt4.SantaFromSecret(fh)
		if err != nil {
			return fmt.Errorf(
				"unable to create santa from secret file %s: %w",
				secretKeyFile,
				err,
			)
		}

		// Add participants
		for _, participantFile := range args {
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
	},
}
