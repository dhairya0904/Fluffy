package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// This command will return if there is a fluffy process runinng behind
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Return the status of fluffy",
	Long:  `Returns if there is a fluffy process running behind`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/fluffy")
		viper.ReadInConfig()

		if viper.GetInt("pID") == -1 {
			fmt.Println("Fluffy not running")
		} else {
			fmt.Printf("Fluffy running with PID: %d\n", viper.GetInt("pID"))
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
