package main

import (
	"fmt"
	"os"

	"github.com/Flontistacks/adl/internal/aria2"
	"github.com/Flontistacks/adl/internal/config"
	"github.com/Flontistacks/adl/internal/tui"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "adl",
		Short: "Terminal download manager powered by aria2c",
		Long: `adl is a Mole-style terminal UI for managing downloads via aria2c.

Run without arguments to open the main menu, or use subcommands for direct access.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI(tui.ViewMenu)
		},
	}

	root.AddCommand(&cobra.Command{
		Use:   "download",
		Short: "Start a new download",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI(tui.ViewDownload)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Show active downloads",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI(tui.ViewList)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "settings",
		Short: "Open settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI(tui.ViewSettings)
		},
	})

	return root
}

func runTUI(view tui.StartView) error {
	if err := aria2.CheckInstalled(); err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if cfg.Aria2Path == "" {
		return fmt.Errorf("aria2c not found. Install with: brew install aria2")
	}
	return tui.Run(cfg, view)
}
