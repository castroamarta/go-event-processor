package cmd

import (
	"github.com/spf13/cobra"
)

func NewProcessCmd() *cobra.Command{
	var file_name string
	var window_size int
	processCmd := &cobra.Command{
		Use:   "go-event-processor --file_name <file-name> --window_size <time-window-in-minutes>",
		// Short: "Event processing command.",
		// Long: `This command receives an event json file as input and 
		// calculates the moving average of the translation delivery time for the last X minutes.`,
	}
	processCmd.Flags().StringVarP(&file_name, "file_name", "f", "", "events input file name")
	processCmd.Flags().IntVarP(&window_size, "window_size", "w", 0, "time window in minutes")
	processCmd.MarkFlagRequired("file_name")
	processCmd.MarkFlagRequired("window_size")
	processCmd.MarkFlagsRequiredTogether("file_name", "window_size")
	return processCmd
}
