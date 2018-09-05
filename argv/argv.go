package argv

import (
	"errors"
	"fmt"
	"strings"
)

// Argv represents a parsed argument with main name
// a list of ops with a `--` prefix and pairs of key
// values.
type Argv struct {
	Name  string
	Sub   *Argv
	Text  string
	Pairs map[string][]string
}

// New returns a new instance of Argv.
func New(name string) *Argv {
	return &Argv{Name: name}
}

// HasKV returns true/false if giving key exists.
func (a *Argv) HasKV(n string) bool {
	_, ok := a.Pairs[n]
	return ok
}

// IsArg returns true/false if giving
// arg has either flags, pairs and a name.
func (a *Argv) IsArg() bool {
	if len(a.Pairs) != 0 {
		return true
	}
	return false
}

// Parse takes provided string, splits according to space
// and parses arguments.
func Parse(args string) (Argv, error) {
	if len(args) == 0 {
		return Argv{}, errors.New("no argument provided")
	}
	return parseArgs(strings.Split(args, " "))
}

// parseArgs attempts to parse the slice of strings
// as a instance of Argv returning an error if one exists.
func parseArgs(args []string) (Argv, error) {
	var argd Argv
	argd.Pairs = map[string][]string{}

	var withCommand bool

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if isIgnored(arg) {
			continue
		}

		// If this is not a flag and are still yet to encounter
		// command then set as command skip.
		if !isFlag(arg) && !withCommand {
			if len(argd.Pairs) != 0 {
				return argd, errors.New("flags must come after command name")
			}

			argd.Name = arg
			withCommand = true
			continue
		}

		// If this is not a flag and we already have a command, possibly
		// a sub-command, to branch out.
		if !isFlag(arg) && withCommand {
			rem := args[i:]
			if len(rem) == 1 {
				if !isFlag(rem[0]) {
					argd.Text = rem[0]
					return argd, nil
				}
			}

			sub, err := parseArgs(rem)
			if err != nil {
				return argd, err
			}

			argd.Sub = &sub
			argd.Text = strings.Join(args[i:], " ")
			return argd, nil
		}

		// check for an equals sign, as in "--foo=bar"
		var key, value string
		var hasEq bool

		opt := strings.TrimLeft(arg, "-")
		if pos := strings.Index(opt, "="); pos != -1 {
			hasEq = true
			key = strings.TrimSpace(opt[:pos])
			value = strings.TrimSpace(opt[pos+1:])
		}

		lastIndex := i
		var values []string

		// deal with the case of list values with either
		// space seperated items or comma seperated items.
		if isList(value) {
			if isListEnd(value) {
				items := strings.TrimSpace(value)
				items = strings.TrimLeft(items, "[")
				items = strings.TrimRight(items, "]")
				values = strings.Split(items, ",")
			} else {
				list := make([]string, 0, 5)
				if before := strings.TrimSpace(strings.TrimLeft(value, "[")); before != "" {
					list = append(list, before)
				}

				for i+1 < len(args) && !isFlag(args[i+1]) && !isListEnd(args[i+1]) {
					list = append(list, strings.TrimSpace(args[i+1]))
					i++
				}

				if isListEnd(args[i+1]) {
					end := strings.TrimSpace(strings.TrimSuffix(args[i+1], "]"))
					list = append(list, end)
					i++
				}

				items := strings.Join(list, " ")
				items = strings.TrimSpace(items)
				items = strings.TrimPrefix(items, "[")
				items = strings.TrimSuffix(items, "]")
				values = strings.Split(items, " ")
			}
		} else if value != "" {
			values = append(values, value)
		}

		// if we have a flag and a equal (=) sign, then it
		// expects a value.
		if key != "" && hasEq && len(values) == 0 {
			return argd, fmt.Errorf("flag %q has no provided value", opt)
		}

		// If we have a flag and there was an eq sign, then its a multi value
		// type as we treat .
		if key != "" && hasEq {
			argd.Pairs[key] = values
			continue
		}

		// if there is a flag and no equal sign existed,  then we probably
		// a branched in sub command, so get last index point, branch out
		// after saving flag into current parent command.
		if opt != "" && key == "" && !hasEq {
			argd.Pairs[opt] = []string{"true"}

			// if we stopped around same index, then
			// push forward.
			if lastIndex != i {
				i = lastIndex
			}
			continue
		}

		if key == "" && !hasEq {
			rem := args[lastIndex:]
			if len(rem) == 1 {
				if !isFlag(rem[0]) {
					argd.Text = rem[0]
					return argd, nil
				}
			}

			sub, err := parseArgs(rem)
			if err != nil {
				return argd, err
			}

			argd.Sub = &sub
			argd.Text = strings.Join(args[i:], " ")
			return argd, nil
		}
	}

	return argd, nil
}

// isFlag returns true if a token is a flag such as "-v" or "--user" but not "-" or "--"
func isFlag(s string) bool {
	return strings.HasPrefix(s, "-") && strings.TrimLeft(s, "-") != ""
}

func isIgnored(s string) bool {
	switch s {
	case "":
		return true
	case "-":
		return true
	case "--":
		return true
	}
	return false
}

func isList(s string) bool {
	return strings.HasPrefix(s, "[")
}

func isListEnd(s string) bool {
	return strings.HasSuffix(s, "]")
}
