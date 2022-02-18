package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/config"
	"github.com/spf13/cobra"
)

func newRunCmd(ctx context.Context, logger *log.Logger, config *config.RdmConfig) *cobra.Command {
	// TODO this needs to diverge, server should hold all the commands and
	// client should query for available commands
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs a custom command defined in the rdm config",
		Run: func(cmd *cobra.Command, args []string) {
			c := client.New()
			content, err := c.SendCommand(ctx, "run", args...)

			if err != nil {
				fmt.Printf("Could not run command: %v", err)
				return
			}

			fmt.Println(string(content))
		},
	}

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		fmt.Println(longRunDescription(config))
	})

	return cmd
}

func longRunDescription(config *config.RdmConfig) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var out strings.Builder

	out.WriteString("Runs a custom command defined via rdm config\n\n")

	c := client.New()
	result, err := c.SendCommand(ctx, "commands")

	if err != nil {
		cancel()
		out.WriteString("Could not communicate with server to get commands")
		return out.String()
	}

	out.WriteString("Available commands:\n")

	for _, command := range bytes.Split(result, []byte("\n")) {
		out.WriteString(fmt.Sprintf("  %s\n", command))
	}

	return out.String()
}
