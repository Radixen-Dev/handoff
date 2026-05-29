package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Radixen-Dev/handoff/internal/config"
	"github.com/Radixen-Dev/handoff/internal/db"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available knowledge packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")

		dbPath, err := config.DBPath()
		if err != nil {
			return err
		}
		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		packages, err := d.List(project)
		if err != nil {
			return err
		}

		if len(packages) == 0 {
			fmt.Fprintln(os.Stderr, "No packages found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROJECT\tTAGS\tEXPIRES")
		for _, pkg := range packages {
			tags := strings.Join(pkg.Tags, ",")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				pkg.ID, pkg.Name, pkg.Project, tags, pkg.ExpiresAt.Format("2006-01-02 15:04"))
		}
		w.Flush()
		return nil
	},
}

func init() {
	listCmd.Flags().String("project", "", "Filter by project")
	rootCmd.AddCommand(listCmd)
}
