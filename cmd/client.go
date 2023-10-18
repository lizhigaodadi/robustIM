package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: nil,
}

func ClientHandle(cmd *cobra.Command, args []string) {
	viper.AddConfigPath()
}
