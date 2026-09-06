package cmd

import (
	"multiviewer-sync/internal/cli/sync"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multiviewer-sync",
	Short: "Syncs together multiple instances of MultiViewer",
	Long: `Allows for a MultiViewer client to be synced to a main instance. This will cause players to play/pause
along with the main instance and the timing to remain synced.`,
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enables verbose output")

	rootCmd.AddCommand(sync.NewCommand())
}
