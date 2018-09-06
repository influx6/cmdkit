package cmdkit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/gokit/cmdkit/argv"
)

const (
	usageTml = `Usage: {{ toLower .Title}} [flags] [command] 

⡿ COMMANDS:{{ range .Commands }}

	⠙ {{toLower .Name }}        {{if isEmpty .ShortDesc }}{{cutoff .Desc 100 }}{{else}}{{cutoff .ShortDesc 100 }}{{end}}
{{end}}
⡿ HELP:

	Run [command] --help to print this message
	Run {{toLower .Title}} --flags to print all flags of all commands.

⡿ Flags:
	{{ range $_, $fl := .Flags }}
	⠙ {{toLower $fl.FlagName}}      {{ if .Default }} Default: {{.Default}} {{end}}     {{ if .Desc }} Desc: {{.Desc}} {{end}}
	{{end}}

`
	flagUsageTml = `Command: {{ toLower .Cmd.Name}} 

⡿ Flags:
	{{$title := toLower .Title}}{{$cmdName := .Cmd.Name}}{{ range $_, $fl := .Cmd.Flags }}
	⠙ {{toLower $fl.FlagName}}      {{ if .Default }} Default: {{.Default}} {{end}}     {{ if .Desc }} Desc: {{.Desc}} {{end}}
	{{end}}
`

	flagOnlyUsageTml = `Usage: {{ toLower .Title}}

⡿ Flags:
	{{ range $_, $fl := .Flags }}
	⠙ {{toLower $fl.FlagName}}      {{ if .Default }} Default: {{.Default}} {{end}}     {{ if .Desc }} Desc: {{.Desc}} {{end}}
	{{end}}
`

	cmdUsageTml = `Command: {{toLower .Cmd.Name}} [flags] [sub commands]

⡿ DESC:

	{{.Cmd.Desc}}

⡿ HELP:

	Run {{toLower .Cmd.Name}} --help to print this message.
	Run {{toLower .Cmd.Name}} [command] --help to print help for sub command.

⡿ Flags:
	{{$title := toLower .Title}}{{$cmdName := .Cmd.Name}}{{ range $_, $fl := .Cmd.Flags }}
	⠙ {{toLower $fl.FlagName}}      {{ if .Default }} Default: {{.Default}} {{end}}     {{ if .Desc }} Desc: {{.Desc}} {{end}}
	{{end}}
⡿ Examples:
	{{ range $_, $content := .Cmd.Usages }}
	⠙ {{$content}}
	{{end}}
⡿ USAGE:
	{{ range $_, $fl := .Cmd.Flags }}
	⠙ {{$title}} --{{toLower $fl.FlagName}}={{.Default}} {{toLower $cmdName}} 
	{{end}}
⡿ SUB COMMANDS:{{ range .Commands }}

	⠙ {{toLower .Name }}        {{if isEmpty .ShortDesc }}{{cutoff .Desc 100 }}{{else}}{{cutoff .ShortDesc 100 }}{{end}}
{{end}}

`
)

var (
	printFlag   = BoolFlag(FlagName("flags"))
	helpFlag    = BoolFlag(FlagName("help"), FlagAlias("h"))
	timeoutFlag = DurationFlag(FlagName("timeout"), FlagAlias("tm"))

	defs = template.FuncMap{
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
		"isEmpty": func(val string) bool {
			return strings.TrimSpace(val) == ""
		},
		"cutoff": func(val string, limit int) string {
			if len(val) > limit {
				return val[:limit] + "..."
			}
			return val
		},
	}
)

// FlagType defines a int to represent a giving flag type.
type FlagType int

// lists of flag types.
const (
	Int FlagType = iota + 1
	Int8
	Int32
	Int16
	Int64
	Bool
	TBool
	String
	Float32
	Float64
	Duration
)

// ValueValidation defines a function type for the purpose
// of validating a giving string input.
type ValueValidation func(string, ...string) error

// ParseFunction defines a function type which is called
// for processing a string.
type ParseFunction func(string, ...string) (interface{}, error)

// MorphFunction defines a function type which takes a
// value and transforms it into another value.
type MorphFunction func(interface{}) (interface{}, error)

// FlagOption defines a function type which takes a giving flagimpl.
type FlagOption func(*Flag)

// Flag implements a structure for parsing string flags.
type Flag struct {
	Name       string
	Alias      string
	Env        string
	Desc       string
	Default    interface{}
	Morph      MorphFunction
	Parser     ParseFunction
	Validation ValueValidation
}

// Validate returns a FlagOption that sets the ValueValidation function.
func Validate(n ValueValidation) FlagOption {
	return func(fl *Flag) {
		fl.Validation = n
	}
}

// Morph sets giving MorphFunction for giving Flag.
func Morph(n MorphFunction) FlagOption {
	return func(fl *Flag) {
		fl.Morph = n
	}
}

// Default returns a FlagOption that sets the desc of a Flag.
func Default(n interface{}) FlagOption {
	return func(fl *Flag) {
		fl.Default = n
	}
}

// FlagDesc returns a FlagOption that sets the desc of a Flag.
func FlagDesc(s string) FlagOption {
	return func(fl *Flag) {
		fl.Desc = s
	}
}

// FlagAlias returns a FlagOption that sets the alias of a Flag.
func FlagAlias(s string) FlagOption {
	return func(fl *Flag) {
		fl.Alias = s
	}
}

// FlagName returns a FlagOption that sets the name of a Flag.
func FlagName(s string) FlagOption {
	return func(fl *Flag) {
		fl.Name = s
	}
}

// Env provides a means to setting the environment variable name
// for a Flag.
func Env(s string) FlagOption {
	return func(fl *Flag) {
		fl.Env = s
	}
}

// FlagAlias returns alias of flag.
func (s *Flag) FlagAlias() string {
	return s.Alias
}

// FlagName returns name of flag.
func (s *Flag) FlagName() string {
	return s.Name
}

// DefaultValue returns Default value of flag pointer.
func (s *Flag) DefaultValue() interface{} {
	return s.Default
}

// Parse sets the underline flag ready for value receiving.
func (s *Flag) Parse(m string, rest ...string) (interface{}, error) {
	if s.Validation != nil {
		if err := s.Validation(m, rest...); err != nil {
			return nil, err
		}
	}

	value, err := s.Parser(m, rest...)
	if err != nil {
		return nil, err
	}

	if s.Morph == nil {
		return value, nil
	}

	return s.Morph(value)
}

// Flags returns the passed in set of variadic arguments
// returning them as a slice.
func Flags(flags ...Flag) []Flag {
	return flags
}

// MakeFlag creates a flag for list of list strings.
func MakeFlag(ops ...FlagOption) Flag {
	var impl Flag
	for _, op := range ops {
		op(&impl)
	}
	return impl
}

// Int64ListFlag creates a flag for list of int64.
func Int64ListFlag(ops ...FlagOption) Flag {
	impl := MakeFlag(ops...)
	if impl.Default != nil {
		if _, ok := impl.Default.([]int64); !ok {
			log.Fatalf("Flag %q must use type []int64 default value types", impl.Name)
		}
	}
	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		if impl.Validation != nil {
			if err := impl.Validation(s, rem...); err != nil {
				return nil, err
			}
		}

		initial, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		elem := make([]int64, 0, 1+len(rem))
		elem = append(elem, initial)

		for _, item := range rem {
			conv, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				return nil, err
			}
			elem = append(elem, conv)
		}
		return elem, nil
	}
	return impl
}

// Float64ListFlag creates a flag for list of list float64.
func Float64ListFlag(ops ...FlagOption) Flag {
	impl := MakeFlag(ops...)
	if impl.Default != nil {
		if _, ok := impl.Default.([]float64); !ok {
			log.Fatalf("Flag %q must use type []float64 default value types", impl.Name)
		}
	}
	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		if impl.Validation != nil {
			if err := impl.Validation(s, rem...); err != nil {
				return nil, err
			}
		}

		initial, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}

		elem := make([]float64, 0, 1+len(rem))
		elem = append(elem, initial)

		for _, item := range rem {
			conv, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return nil, err
			}
			elem = append(elem, conv)
		}
		return elem, nil
	}
	return impl
}

// BoolListFlag creates a flag for list of bool.
func BoolListFlag(ops ...FlagOption) Flag {
	impl := MakeFlag(ops...)
	if impl.Default != nil {
		if _, ok := impl.Default.([]bool); !ok {
			log.Fatalf("Flag %q must use type []bool default value types", impl.Name)
		}
	}
	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		if impl.Validation != nil {
			if err := impl.Validation(s, rem...); err != nil {
				return nil, err
			}
		}

		initial, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}

		elem := make([]bool, 0, 1+len(rem))
		elem = append(elem, initial)

		for _, item := range rem {
			conv, err := strconv.ParseBool(item)
			if err != nil {
				return nil, err
			}
			elem = append(elem, conv)
		}
		return elem, nil
	}
	return impl
}

// UIntListFlag creates a flag for list of uint.
func UIntListFlag(ops ...FlagOption) Flag {
	impl := MakeFlag(ops...)
	if impl.Default != nil {
		if _, ok := impl.Default.([]uint); !ok {
			log.Fatalf("Flag %q must use type []uint default value types", impl.Name)
		}
	}
	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		if impl.Validation != nil {
			if err := impl.Validation(s, rem...); err != nil {
				return nil, err
			}
		}

		initial, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		elem := make([]uint, 0, 1+len(rem))
		elem = append(elem, uint(initial))

		for _, item := range rem {
			conv, err := strconv.ParseUint(item, 10, 64)
			if err != nil {
				return nil, err
			}
			elem = append(elem, uint(conv))
		}
		return elem, nil
	}
	return impl
}

// IntListFlag creates a flag for list of int.
func IntListFlag(ops ...FlagOption) Flag {
	impl := MakeFlag(ops...)
	if impl.Default != nil {
		if _, ok := impl.Default.([]int); !ok {
			log.Fatalf("Flag %q must use type []int default value types", impl.Name)
		}
	}
	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		if impl.Validation != nil {
			if err := impl.Validation(s, rem...); err != nil {
				return nil, err
			}
		}

		initial, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		elem := make([]int, 0, 1+len(rem))
		elem = append(elem, int(initial))

		for _, item := range rem {
			conv, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				return nil, err
			}
			elem = append(elem, int(conv))
		}
		return elem, nil
	}
	return impl
}

// StringListFlag creates a flag for list of list strings.
func StringListFlag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.([]string); !ok {
			log.Fatalf("Flag %q must use type []string default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		if impl.Validation != nil {
			if err := impl.Validation(s, rem...); err != nil {
				return nil, err
			}
		}

		return append([]string{s}, rem...), nil
	}
	return impl
}

// StringFlag creates a flag for strings.
func StringFlag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(string); !ok {
			log.Fatalf("Flag %q must use type string default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		return s, nil
	}
	return impl
}

// TBoolFlag creates a flag for duration.
func TBoolFlag(ops ...FlagOption) Flag {
	var impl Flag
	impl.Default = true
	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(bool); !ok {
			log.Fatalf("Flag %q must use type bool default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseBool(s)
		if err != nil {
			return nil, errors.New("not a bool")
		}
		return myValue, nil
	}
	return impl
}

// BoolFlag creates a flag for duration.
func BoolFlag(ops ...FlagOption) Flag {
	var impl Flag

	impl.Default = false
	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(bool); !ok {
			log.Fatalf("Flag %q must use type bool default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseBool(s)
		if err != nil {
			return nil, errors.New("not a bool")
		}
		return myValue, nil
	}
	return impl
}

// DurationFlag creates a flag for duration.
func DurationFlag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(time.Duration); !ok {
			log.Fatalf("Flag %q must use type time.Duration default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := time.ParseDuration(s)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// Int8Flag creates a flag for int8.
func Int8Flag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(int8); !ok {
			log.Fatalf("Flag %q must use type int8 default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// Int16Flag creates a flag for int16.
func Int16Flag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(int16); !ok {
			log.Fatalf("Flag %q must use type int16 default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// IntFlag creates a flag for int.
func IntFlag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(int); !ok {
			log.Fatalf("Flag %q must use type int default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// Float64Flag creates a flag for int.
func Float64Flag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(float64); !ok {
			log.Fatalf("Flag %q must use type float64 default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// Float32Flag creates a flag for int.
func Float32Flag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(float32); !ok {
			log.Fatalf("Flag %q must use type float32 default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// Int64Flag creates a flag for int.
func Int64Flag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(int64); !ok {
			log.Fatalf("Flag %q must use type int64 default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// Int32Flag creates a flag for int.
func Int32Flag(ops ...FlagOption) Flag {
	var impl Flag

	for _, op := range ops {
		op(&impl)
	}

	if impl.Default != nil {
		if _, ok := impl.Default.(int32); !ok {
			log.Fatalf("Flag %q must use type int32 default value types", impl.Name)
		}
	}

	impl.Parser = func(s string, rem ...string) (interface{}, error) {
		myValue, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, errors.New("not a int")
		}
		return myValue, nil
	}
	return impl
}

// Action defines a giving function to be executed for a Command.
type Action func(Context) error

// KeyValue defines a interface that expose giving set of methods
// for retrieving key and value from context.
type KeyValue interface {
	IsSet(string) bool
	Int(string) int
	Bool(string) bool
	Uint(string) uint
	Uint64(string) uint64
	Int64(string) int64
	String(string) string
	Float64(string) float64
	Duration(string) time.Duration
	Get(string) (interface{}, bool)
}

// Context defines a interface which combines the bag.ter with a
// provided context.
type Context interface {
	KeyValue

	PrintHelp()
	Args() []string
	Parent() KeyValue
	Ctx() context.Context
}

type ctxImpl struct {
	ctx         context.Context
	args        []string
	HelpPrinter func()
	parent      Context
	flags       map[string]struct{}
	pairs       map[string]interface{}
}

// Args returning the internal associated arg list.
// It implements the Context interface.
func (c ctxImpl) Args() []string {
	return c.args
}

// Ctx returns the context.Context associated with the command context.
func (c ctxImpl) Ctx() context.Context {
	return c.ctx
}

// Parent returns a Context that is the context of
// a parent command in relation to the command that
// generated this context.
func (c ctxImpl) Parent() KeyValue {
	return c.parent
}

// PrintHelp calls underline function to print help for command.
func (c ctxImpl) PrintHelp() {
	if c.HelpPrinter != nil {
		c.HelpPrinter()
	}
}

// Duration returns the duration value of a key if it exists.
func (c *ctxImpl) Duration(key string) time.Duration {
	if val, found := c.Get(key); found {
		return val.(time.Duration)
	}
	return 0
}

// Bool returns the bool value of a key if it exists.
func (c *ctxImpl) Bool(key string) bool {
	if val, found := c.Get(key); found {
		return val.(bool)
	}
	return false
}

// Float64 returns the float64 value of a key if it exists.
func (c *ctxImpl) Float64(key string) float64 {
	if val, found := c.Get(key); found {
		return val.(float64)
	}
	return 0
}

// Float32 returns the float32 value of a key if it exists.
func (c *ctxImpl) Float32(key string) float32 {
	if val, found := c.Get(key); found {
		return val.(float32)
	}
	return 0
}

// Int8 returns the int8 value of a key if it exists.
func (c *ctxImpl) Int8(key string) int8 {
	if val, found := c.Get(key); found {
		return val.(int8)
	}
	return 0
}

// Int16 returns the int16 value of a key if it exists.
func (c *ctxImpl) Int16(key string) int16 {
	if val, found := c.Get(key); found {
		return val.(int16)
	}
	return 0
}

// Int64 returns the value type value of a key if it exists.
func (c *ctxImpl) Int64(key string) int64 {
	if val, found := c.Get(key); found {
		return val.(int64)
	}
	return 0
}

// Int32 returns the value type value of a key if it exists.
func (c *ctxImpl) Int32(key string) int32 {
	if val, found := c.Get(key); found {
		return val.(int32)
	}
	return 0
}

// Uint64 returns the value type value of a key if it exists.
func (c *ctxImpl) Uint64(key string) uint64 {
	if val, found := c.Get(key); found {
		return val.(uint64)
	}
	return 0
}

// Uint returns the value type value of a key if it exists.
func (c *ctxImpl) Uint(key string) uint {
	if val, found := c.Get(key); found {
		return val.(uint)
	}
	return 0
}

// Int returns the value type value of a key if it exists.
func (c *ctxImpl) Int(key string) int {
	if val, found := c.Get(key); found {
		return val.(int)
	}
	return 0
}

// String returns the value type value of a key if it exists.
func (c *ctxImpl) String(key string) string {
	if val, found := c.Get(key); found {
		return val.(string)
	}
	return ""
}

// Get returns the value of a key if it exists.
// If the key is not seen within present context, then the parent
// of context is checked for giving key.
func (c *ctxImpl) Get(key string) (value interface{}, found bool) {
	if item, ok := c.pairs[key]; ok {
		return item, true
	}
	if c.parent == nil {
		return nil, false
	}
	return c.parent.Get(key)
}

// IsSet returns true/false if giving key was set in command context.
func (c *ctxImpl) IsSet(key string) bool {
	if _, ok := c.pairs[key]; ok {
		return true
	}
	return false
}

func (c *ctxImpl) process(arg *argv.Argv, flags []Flag) error {
	if c.pairs == nil {
		c.flags = map[string]struct{}{}
		c.pairs = map[string]interface{}{}
	}

	for _, flag := range flags {
		c.flags[flag.FlagName()] = struct{}{}
		c.flags[flag.FlagAlias()] = struct{}{}
		if flagValue, provided := arg.Pairs[flag.FlagName()]; provided {
			value, err := flag.Parse(flagValue[0], flagValue[1:]...)
			if err != nil {
				return err
			}
			c.pairs[flag.FlagName()] = value
			c.pairs[flag.FlagAlias()] = value
			continue
		}
		if flag.Env != "" {
			value, err := flag.Parse(os.Getenv(flag.Env))
			if err != nil {
				return err
			}
			c.pairs[flag.FlagName()] = value
			c.pairs[flag.FlagAlias()] = value
			continue
		}
		if flag.DefaultValue() != nil {
			c.pairs[flag.FlagName()] = flag.DefaultValue()
			c.pairs[flag.FlagAlias()] = flag.DefaultValue()
		}
	}
	return nil
}

// CommandFunc defines a function type that modifies a giving Command.
type CommandFunc func(*Command)

// ShortDesc sets giving name for provided command.
func ShortDesc(desc string) CommandFunc {
	return func(cmd *Command) {
		cmd.ShortDesc = desc
	}
}

// Desc sets giving name for provided command.
func Desc(desc string) CommandFunc {
	return func(cmd *Command) {
		cmd.Desc = desc
	}
}

// WithAction sets giving name for provided command.
func WithAction(ac Action) CommandFunc {
	return func(cmd *Command) {
		cmd.Action = ac
	}
}

// Usage sets adds usage text for provided command.
func Usage(desc string) CommandFunc {
	return func(cmd *Command) {
		cmd.Usages = append(cmd.Usages, desc)
	}
}

// SubCommands adds giving commands into command list of
// parent.
func SubCommands(cms ...Command) CommandFunc {
	return func(cmd *Command) {
		for _, cm := range cms {
			cmd.Commands[cm.Name] = cm
		}
	}
}

// Command defines structures which define specific actions to be executed
// with associated flags.
// Commands provided will have their ShortDesc trimmed to 100 in length, so
// ensure to have what you wanna say fit 100 and put more detail explanations
// in Desc field.
type Command struct {
	Name         string
	Desc         string
	ShortDesc    string
	Action       Action
	Flags        []Flag
	Usages       []string
	FlagUsage    string
	CommandUsage string
	Stderr       io.Writer
	Stdout       io.Writer
	Commands     map[string]Command
}

// Run executes giving command with argv.Argv object.
func (c *Command) Run(arg *argv.Argv, parent Context) error {
	if arg.HasKV("help") || arg.HasKV("h") {
		_, err := fmt.Fprint(c.Stderr, c.CommandUsage)
		return err
	}

	if arg.HasKV("flags") {
		_, err := fmt.Fprint(c.Stderr, c.FlagUsage)
		return err
	}

	if c.Action == nil {
		return fmt.Errorf("no action associated with command %q", c.Name)
	}

	var childCtx ctxImpl
	childCtx.parent = parent
	childCtx.ctx = parent.Ctx()
	if err := childCtx.process(arg, c.Flags); err != nil {
		return err
	}

	// if we are dealing with possible tree then go down the tree.
	if arg.Sub != nil {
		return c.runSubCommand(arg.Sub, &childCtx)
	}

	// if we are dealing with the final argv, then is the it's text
	// value a command also, if it is, make a new chain and pass it on.
	if _, ok := c.Commands[arg.Text]; ok {
		arg.Sub = argv.New(arg.Text)
		return c.runSubCommand(arg.Sub, &childCtx)
	}

	cancel := func() {}
	ctx := parent.Ctx()
	if tm := childCtx.Duration("timeout"); childCtx.IsSet("timeout") {
		childCtx.ctx, cancel = context.WithTimeout(ctx, tm)
	}

	defer cancel()

	return c.Action(&childCtx)
}

func (c *Command) runSubCommand(arg *argv.Argv, parent Context) error {
	for _, sub := range c.Commands {
		if sub.Name == arg.Name {
			return sub.Run(arg, parent)
		}
	}
	return fmt.Errorf("%q has no subcommand named %q", c.Name, arg.Name)
}

// Commands returns the passed in set of variadic arguments
// returning them as a slice.
func Commands(cmds ...Command) []Command {
	return cmds
}

// Cmd returns a new Command from the provided options.
func Cmd(name string, ops ...CommandFunc) Command {
	cm := Command{
		Stderr:   os.Stderr,
		Stdout:   os.Stdout,
		Commands: map[string]Command{},
		Name:     strings.ToLower(name),
	}

	for _, op := range ops {
		op(&cm)
	}

	if tml, err := template.New("command.Usage").Funcs(defs).Parse(cmdUsageTml); err == nil {
		var bu bytes.Buffer
		if err := tml.Execute(&bu, struct {
			Title    string
			Cmd      Command
			Commands map[string]Command
		}{
			Cmd:      cm,
			Title:    cm.Name,
			Commands: cm.Commands,
		}); err != nil {
			log.Fatalf("Error occured compiling command %q usage text: %q", cm.Name, err)
		}

		cm.CommandUsage = bu.String()
	}

	if tml, err := template.New("flags.Usage").Funcs(defs).Parse(flagUsageTml); err == nil {
		var bu bytes.Buffer
		if err := tml.Execute(&bu, struct {
			Title    string
			Cmd      Command
			Commands map[string]Command
		}{
			Cmd:      cm,
			Title:    cm.Name,
			Commands: cm.Commands,
		}); err != nil {
			log.Fatalf("Error occured compiling command %q flag usage text: %q", cm.Name, err)
		}

		cm.FlagUsage = bu.String()
	}

	return cm
}

// Run adds all commands and appropriate flags for each commands.
// There is no need to call flag.Parse, has this calls it underneath and
// parses appropriate commands.
func Run(title string, flags []Flag, cmds []Command) {
	title = strings.ToLower(title)
	commands := map[string]Command{}

	flags = append(flags, helpFlag)
	flags = append(flags, printFlag)
	flags = append(flags, timeoutFlag)

	// Register all flags first.
	for _, cmd := range cmds {
		commands[cmd.Name] = cmd
	}

	var cmdHelp string
	var flagHelp string

	tml, err := template.New("command.Usage").Funcs(defs).Parse(usageTml)
	if err != nil {
		log.Fatal("Failed to create template instance: ", err)
	}

	tmlflags, err := template.New("flags.Usage").Funcs(defs).Parse(flagOnlyUsageTml)
	if err != nil {
		log.Fatal("Failed to create template instance: ", err)
	}

	var bu bytes.Buffer
	if err = tml.Execute(&bu, struct {
		Title    string
		Commands []Command
		Flags    []Flag
	}{
		Title:    title,
		Flags:    flags,
		Commands: cmds,
	}); err != nil {
		log.Fatal("Failed to generated help message for command: ", err)
	}
	cmdHelp = bu.String()

	bu.Reset()
	if err = tmlflags.Execute(&bu, struct {
		Title string
		Flags []Flag
	}{
		Title: title,
		Flags: flags,
	}); err != nil {
		log.Fatal("Failed to generated help message for command: ", err)
	}
	flagHelp = bu.String()

	args := strings.Join(os.Args, " ")
	carg, err := argv.Parse(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	// if we are dealing with the final argv, then is the it's text
	// value a command also, if it is, make a new chain and pass it on.
	if _, ok := commands[carg.Text]; ok {
		carg.Sub = argv.New(carg.Text)
	}

	if carg.HasKV("h") || carg.HasKV("help") {
		fmt.Fprint(os.Stderr, cmdHelp)
		return
	}

	if carg.HasKV("flags") {
		fmt.Fprint(os.Stderr, flagHelp)
		return
	}

	if carg.Sub == nil {
		fmt.Fprint(os.Stderr, cmdHelp)
		return
	}

	target, ok := commands[carg.Sub.Name]
	if !ok {
		fmt.Fprint(os.Stderr, fmt.Errorf("command not found %q", carg.Name))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cmdCtx ctxImpl
	cmdCtx.ctx = ctx
	if err := cmdCtx.process(&carg, flags); err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGQUIT)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		defer close(ch)
		if err := target.Run(carg.Sub, &cmdCtx); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}
	}()

	<-ch
}
