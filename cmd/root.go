package cmd

import (
	"os"

	"github.com/I1820/gm/cmd/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ExitFailure status code
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var root = &cobra.Command{
		Use:   "gm",
		Short: "Who decrypts the LoRa payloads",
	}

	root.Println("13 Feb 2020, Best Day Ever")

	server.Register(root)

	if err := root.Execute(); err != nil {
		logrus.Errorf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
