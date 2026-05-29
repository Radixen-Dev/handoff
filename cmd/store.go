package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Radixen-Dev/handoff/internal/config"
	"github.com/Radixen-Dev/handoff/internal/db"
	"github.com/Radixen-Dev/handoff/internal/model"

	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store a knowledge package (reads content from stdin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		summary, _ := cmd.Flags().GetString("summary")
		ttl, _ := cmd.Flags().GetString("ttl")
		project, _ := cmd.Flags().GetString("project")
		tagsStr, _ := cmd.Flags().GetString("tags")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		// Parse TTL
		duration, err := parseTTL(ttl)
		if err != nil {
			return fmt.Errorf("invalid --ttl: %w", err)
		}

		// Read content from stdin
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		if len(strings.TrimSpace(string(content))) == 0 {
			return fmt.Errorf("no content provided on stdin")
		}

		// Parse tags
		var tags []string
		if tagsStr != "" {
			for _, t := range strings.Split(tagsStr, ",") {
				tags = append(tags, strings.TrimSpace(t))
			}
		}

		// Generate ID
		id := generateID()

		now := time.Now().UTC()
		pkg := &model.Package{
			ID:        id,
			Name:      name,
			Summary:   summary,
			Content:   string(content),
			Tags:      tags,
			Project:   project,
			CreatedAt: now,
			ExpiresAt: now.Add(duration),
		}

		// Open DB and store
		dbPath, err := config.DBPath()
		if err != nil {
			return err
		}
		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		if err := d.Store(pkg); err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "%s\n", id)
		return nil
	},
}

func init() {
	storeCmd.Flags().String("name", "", "Name for the knowledge package (required)")
	storeCmd.Flags().String("summary", "", "Short summary of the package contents")
	storeCmd.Flags().String("ttl", "7d", "Time-to-live (e.g., 2h, 7d, 30d)")
	storeCmd.Flags().String("project", "", "Project grouping key")
	storeCmd.Flags().String("tags", "", "Comma-separated tags")
	rootCmd.AddCommand(storeCmd)
}

func parseTTL(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("too short")
	}

	unit := s[len(s)-1]
	val := s[:len(s)-1]

	var multiplier time.Duration
	switch unit {
	case 'h':
		multiplier = time.Hour
	case 'd':
		multiplier = 24 * time.Hour
	default:
		return 0, fmt.Errorf("unsupported unit '%c' (use h or d)", unit)
	}

	var n int
	if _, err := fmt.Sscanf(val, "%d", &n); err != nil {
		return 0, fmt.Errorf("invalid number: %s", val)
	}
	if n <= 0 {
		return 0, fmt.Errorf("must be positive")
	}

	return time.Duration(n) * multiplier, nil
}

func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
