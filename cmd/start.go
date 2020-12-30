package cmd

import (
	"context"
	"fluffy/alert"
	"fluffy/domain"
	"fluffy/monitor"
	"fluffy/reader"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	REPORT_TIME  = "report-time"
	ALERT_WINDOW = "alert-window"
	THRESHOLD    = "threshold"
	FILENAME     = "fileName"
	REPORT_PATH  = "report-path"
	LOG_PATH     = "log-path"
	ALERT_PATH   = "alert-path"
)

// this command will start log monitoring, config will be picked from config.yml
// Run this using shell script start-server to run as a daemon
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Fluffy will start monitoring your logs",
	Long:  `Set configuration in config.yml. Fluffy will use the config and populate all reports and alerts in given file`,
	Run: func(cmd *cobra.Command, args []string) {

		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.ReadInConfig()
		reportTime := viper.GetInt(REPORT_TIME)
		alertWindow := viper.GetInt(ALERT_WINDOW)
		threshold := viper.GetInt(THRESHOLD)
		fileName := viper.GetString(FILENAME)
		reportPath := viper.GetString(REPORT_PATH)
		logPath := viper.GetString(LOG_PATH)
		alertPath := viper.GetString(ALERT_PATH)

		if reportTime == 0 || alertWindow == 0 || threshold == 0 || fileName == "" {
			fmt.Println("config.yml: bad configuration")
			os.Exit(1)
		}

		startMonitoring(reportTime, alertWindow, threshold, fileName, reportPath, logPath, alertPath)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

/// Will start monitoring the given log file.
func startMonitoring(reportTime, alertWindow, threshold int, fileName, reportPath, logPath, alertPath string) {

	pID := os.Getpid()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.Set("pID", pID)
	viper.WriteConfig()

	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		os.Exit(1)
	}
	defer f.Close()
	log.SetOutput(f)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(3)

	reader := reader.NewReader()
	eventsMonitor := make(chan domain.Event)
	eventsAlert := make(chan domain.Event)

	reader.Subscribe(eventsMonitor)
	reader.Subscribe(eventsAlert)

	monitor := monitor.NewMonitor(reportTime)
	alertMonitor := alert.NewAlertMonitor(alertWindow, threshold)

	go reader.StartPublishing(ctx, &wg, fileName)
	go monitor.StartMonitor(ctx, &wg, eventsMonitor, reportPath)
	go alertMonitor.StartAlertMonitor(ctx, &wg, eventsAlert, alertPath)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println()
		log.Println(sig)
		cancel()
	}()

	<-ctx.Done()

	viper.Set("pID", -1)
	viper.WriteConfig()

	wg.Wait()
}
