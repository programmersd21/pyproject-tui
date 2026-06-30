// Package main provides the pyproject-tui command-line entrypoint.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/programmersd21/pyproject-tui/internal/model"
	"github.com/programmersd21/pyproject-tui/internal/parser"
	"github.com/spf13/cobra"
)

var (
	version = versionFromFile()
	commit  = "none"
	date    = "unknown"
)

func main() {
	var create bool
	var showVersion bool
	root := &cobra.Command{
		Use:   "pyproject-tui [path]",
		Short: "Keyboard-driven TUI for pyproject.toml",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if showVersion {
				fmt.Printf("pyproject-tui v%s (commit %s, built %s)\n", version, commit, date)
				return nil
			}
			path := "./pyproject.toml"
			if len(args) == 1 {
				path = args[0]
			}
			if info, statErr := os.Stat(path); statErr != nil {
				if os.IsNotExist(statErr) && create {
					pp := parser.NewEmpty(path)
					if writeErr := parser.Write(pp); writeErr != nil {
						return writeErr
					}
					return runTUI(pp, false)
				}
				if os.IsNotExist(statErr) {
					return fmt.Errorf("pyproject file not found: %s", path)
				}
				return statErr
			} else if info.IsDir() {
				return fmt.Errorf("%s is a directory", path)
			}
			pp, loadErr := parser.Load(path)
			if loadErr != nil {
				if parser.IsSyntaxError(loadErr) {
					pp = parser.NewEmpty(path)
					if raw, rerr := parser.LoadRaw(path); rerr == nil {
						pp.Raw = raw
					}
					return runTUIWithStatus(pp, true, loadErr.Error())
				}
				return loadErr
			}
			return runTUI(pp, false)
		},
	}
	root.Flags().BoolVar(&create, "create", false, "create a default pyproject.toml if missing")
	root.Flags().BoolVarP(&showVersion, "version", "v", false, "print version")
	root.SilenceUsage = true
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runTUI(pp *parser.PyProject, readOnly bool) error {
	app := model.NewAppModel(pp, readOnly)
	app.SetVersion(version)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

func runTUIWithStatus(pp *parser.PyProject, readOnly bool, status string) error {
	app := model.NewAppModel(pp, readOnly)
	app.SetVersion(version)
	app.SetStatus(status, true)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
