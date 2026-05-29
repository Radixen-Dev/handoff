package cmd

import (
	"fmt"

	"github.com/Radixen-Dev/handoff/internal/config"
	"github.com/Radixen-Dev/handoff/internal/db"

	"github.com/spf13/cobra"
)

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "Remove all expired packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPath, err := config.DBPath()
		if err != nil {
			return err
		}
		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		n, err := d.GC()
		if err != nil {
			return err
		}

		fmt.Printf("Removed %d expired package(s).\n", n)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gcCmd)
}
