package cmd

import (
	"fmt"
	"os"

	"github.com/Radixen-Dev/handoff/internal/config"
	"github.com/Radixen-Dev/handoff/internal/db"

	"github.com/spf13/cobra"
)

var retrieveCmd = &cobra.Command{
	Use:   "retrieve [id]",
	Short: "Retrieve a knowledge package by ID or name",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")

		if len(args) == 0 && name == "" {
			return fmt.Errorf("provide a package ID as argument or use --name")
		}

		dbPath, err := config.DBPath()
		if err != nil {
			return err
		}
		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		if len(args) > 0 {
			pkg, err := d.GetByID(args[0])
			if err != nil {
				return err
			}
			fmt.Fprint(os.Stdout, pkg.Content)
			return nil
		}

		pkg, err := d.GetByName(name)
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, pkg.Content)
		return nil
	},
}

func init() {
	retrieveCmd.Flags().String("name", "", "Retrieve by package name (gets most recent)")
	rootCmd.AddCommand(retrieveCmd)
}
