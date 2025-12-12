package cmd

import "strings"

// Command is a simple placeholder command.
type Command struct {
	name string
}

func (c *Command) Name() string { return c.name }

// ParseArgs is a lightweight argument parser compatible with TeaGo/cmd.ParseArgs.
// It splits by space while keeping quoted segments intact.
func ParseArgs(line string) []string {
	fields := []string{}
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(line); i++ {
		ch := line[i]
		switch ch {
		case ' ', '\t':
			if inQuote {
				current.WriteByte(ch)
			} else if current.Len() > 0 {
				fields = append(fields, current.String())
				current.Reset()
			}
		case '\'', '"':
			if inQuote && ch == quoteChar {
				inQuote = false
			} else if !inQuote {
				inQuote = true
				quoteChar = ch
			} else {
				current.WriteByte(ch)
			}
		default:
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		fields = append(fields, current.String())
	}
	return fields
}

// AllCommands placeholder returning empty command map.
func AllCommands() map[string]*Command { return map[string]*Command{} }

// Try placeholder runner.
// In TeaGo, Try returns whether command found; here always false.
func Try(args ...any) bool { return false }
