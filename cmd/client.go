package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: nil,
}

func ClientHandle(cmd *cobra.Command, args []string) {
	/*TODO: Add a command parsing module*/
}
