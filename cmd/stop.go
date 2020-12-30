package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var FLUFFY_NOT_RUNNING = "Fluffy not running"

// Will stop if there is a fluffy process runnning
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running fluffy process",
	Long:  `Stop fluffy if there is a process running`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.ReadInConfig()

		if viper.GetInt("pID") == -1 {
			fmt.Println(FLUFFY_NOT_RUNNING)
		} else {
			fmt.Printf("Fluffy running with PID: %d\n", viper.GetInt("pID"))
			p, err := os.FindProcess(viper.GetInt("pID"))

			if err != nil {
				fmt.Println(FLUFFY_NOT_RUNNING)
			}
			err = p.Signal(os.Interrupt)
			if err != nil {
				fmt.Println(FLUFFY_NOT_RUNNING)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
