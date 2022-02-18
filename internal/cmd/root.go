package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/blakewilliams/remote-development-manager/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rdm",
	Short: "A server and client for better remote development integration.",
}

var LogPath string = os.TempDir() + "rdm.log"

func Execute(ctx context.Context) error {
	logger := log.Default()

	rdmConfig := config.New()

	readConfig(rdmConfig)

	rootCmd.AddCommand(newServerCmd(ctx, logger, rdmConfig))
	rootCmd.AddCommand(newCopyCmd(ctx, logger))
	rootCmd.AddCommand(newPasteCmd(ctx, logger))
	rootCmd.AddCommand(newOpenCmd(ctx, logger))
	rootCmd.AddCommand(newSocketCmd(ctx))
	rootCmd.AddCommand(newStopCmd(ctx, logger))
	rootCmd.AddCommand(newLogpathCmd(ctx))
	rootCmd.AddCommand(newRunCmd(ctx, logger, rdmConfig))

	if rdmConfig != nil {
		rootCmd.AddCommand(newRunCmd(ctx, logger, rdmConfig))
	}

	return rootCmd.Execute()
}

func readConfig(rdmConfig *config.RdmConfig) {
	home, err := os.UserHomeDir()
	path := filepath.Join(home, ".config/rdm/rdm.json")

	_, err = os.Stat(path)

	if err != nil {
		return
	}

	contents, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(contents, rdmConfig)

	if err != nil {
		panic(err)
	}
}
