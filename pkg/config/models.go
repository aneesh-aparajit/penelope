package config

import "fmt"

type CLI int

const (
	Claude CLI = iota
	Codex
	Cursor
)

type Model int

const (
	Opus Model = iota
	Sonnet
	Haiku
)

func (m *Model) UnmarshalText(text []byte) error {
	switch string(text) {
	case "opus":
		*m = Opus
	case "sonnet":
		*m = Sonnet
	case "haiku":
		*m = Haiku
	default:
		return fmt.Errorf("unknown model: %s", string(text))
	}
	return nil
}

func (c *CLI) UnmarshalText(text []byte) error {
	switch string(text) {
	case "claude":
		*c = Claude
	case "codex":
		*c = Codex
	case "cursor":
		*c = Cursor
	default:
		return fmt.Errorf("unknown CLI: %s", string(text))
	}
	return nil
}
