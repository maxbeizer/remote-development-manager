package cmd

import (
	"context"
	"log"
	"os"

	"github.com/blakewilliams/remote-development-manager/internal/clipboard"
	"github.com/blakewilliams/remote-development-manager/internal/config"
	"github.com/blakewilliams/remote-development-manager/internal/server"
	"github.com/spf13/cobra"
)

func newServerCmd(ctx context.Context, logger *log.Logger, rdmConfig *config.RdmConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Starts a server on the local machine.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			logFile, err := os.OpenFile(LogPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
			defer logFile.Close()
			log.SetOutput(logFile)

			s := server.New(server.UnixSocketPath(), clipboard.MacosClipboard, logger, rdmConfig)
			err = s.Listen(ctx)

			if err != nil {
				log.Printf("Server could not be started: %v", err)
				cancel()
				return
			}
		},
	}
}
