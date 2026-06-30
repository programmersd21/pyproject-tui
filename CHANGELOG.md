# Changelog

All notable changes to this project will be documented in this file.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-06-30 (Intial release)
### Added
- Full `pyproject.toml` parsing (all sections)
- Two-pane TUI: section sidebar + field editor
- Inline editing for all field types (string, array, map, table-array)
- Generic tree rendering for `[tool.*]` sections
- Dirty-state tracking and save with `s`
- Help overlay with `?`
- `--version` flag with build metadata
- Semantic theme system with centralized style generation
- Persistent settings stored in the OS config directory
- Live theme switching without restart
- Dashboard, settings, and help views
- Dashboard open-file prompt for jumping directly to a `pyproject.toml`
- Native config open - press `i` to open the config file in the OS default editor
- Add new `[tool.*]` section - press `a` in the sidebar
- Undo/redo with 50-level history - `u` undoes, `r` redoes
- Responsive sidebar width - adapts to terminal size (25/28/32 chars)
- Visible sidebar divider - explicit `│` character between sidebar and editor
- Focus system - `Tab`/`Shift+Tab` switches between sidebar and editor with clear visual feedback (active pane colorful, inactive pane gray)
- Focus indicator in footer - active pane name shown in accent color
- App version in footer - displayed in green at the right edge of the status bar
- Sandfall reveal animation - tiny `·` grains fall from the top in a wavy stream and settle into the PYPROJECT TUI logo
- Smooth gradient animation - after sandfall completes, logo transitions to a sine-wave gradient using the current theme's colors
- Theme-cycled sandfall restart - pressing `t` restarts the sandfall with the new theme palette
- 10 premium themes including Tokyo Night, Catppuccin, Nord, Gruvbox, Rose Pine, Everforest, Python, Midnight, Minimal, Sage
- Settings page with compact layout, Unicode section icons, rounded border, and live theme preview
- Python theme preset with a Python-inspired blue and gold palette
- Immediate settings propagation for theme, density, borders, and line numbers
- Cross-platform Makefile for Windows, Linux, and macOS
- `workflow_dispatch` trigger on CI workflow for manual runs
