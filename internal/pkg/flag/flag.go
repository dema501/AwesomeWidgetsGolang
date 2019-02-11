package flag

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrHelp              = errors.New("flag: help requested")
	ErrFlagNotDefined    = errors.New("flag provided but not defined")
	ErrFlagBadSyntax     = errors.New("bad flag syntax")
	ErrFlagNeedsArgument = errors.New("flag needs an argument")
	ErrFlagInvalidValue  = errors.New("invalid value")
)

// FlagSet represents a set of defined flags. The zero value of a FlagSet
// has no name and has ContinueOnError error handling.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler. What happens after Usage is called depends
	// on the ErrorHandling setting; for the command line, this defaults
	// to ExitOnError, which exits the program after calling Usage.
	Usage  func()
	name   string
	parsed bool
	actual map[string]*Flag
	formal map[string]*Flag
	args   []string // arguments after flags
}

// A Flag represents the state of a flag.
type Flag struct {
	Name      string // name as it appears on command line
	Usage     string // help message
	Transform func(string) error
}

// Val defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (f *FlagSet) Val(name string, usage string, transform func(string) error) {
	flag := &Flag{
		Name:      name,
		Usage:     usage,
		Transform: transform,
	}
	_, alreadythere := f.formal[name]
	if alreadythere {
		var msg string
		if f.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", f.name, name)
		}
		fmt.Println(msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}

	f.formal[name] = flag
}

// parseOne parses one flag. It reports whether a flag was seen.
func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}
	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, errors.Wrapf(ErrFlagBadSyntax, "%s", s)
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}
	m := f.formal
	flag, alreadythere := m[name] // BUG
	if !alreadythere {
		if name == "help" || name == "h" { // special case for nice help message.
			f.usage()
			return false, ErrHelp
		}
		return false, errors.Wrapf(ErrFlagNotDefined, "-%s", name)
	}

	// It must have a value, which might be the next argument.
	if !hasValue && len(f.args) > 0 {
		// value is the next arg
		hasValue = true
		value, f.args = f.args[0], f.args[1:]
	}
	if !hasValue {
		return false, errors.Wrapf(ErrFlagNeedsArgument, "-%s", name)
	}

	if err := flag.Transform(value); err != nil {
		return false, errors.Wrapf(ErrFlagInvalidValue, "%q for flag -%s: %v", value, name, err)
	}

	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flag
	return true, nil
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.
func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}

		return err
	}
	return nil
}

// VisitAll visits the flags in lexicographical order, calling fn for each.
// It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flag := range f.formal {
		fn(flag)
	}
}

func (f *FlagSet) usage() {
	if f.Usage == nil {
		f.DefaultUsage()
	} else {
		f.Usage()
	}
}

// defaultUsage is the default function to print a usage message.
func (f *FlagSet) DefaultUsage() {
	if f.name == "" {
		fmt.Printf("Usage:\n")
	} else {
		fmt.Printf("Usage of %s:\n", f.name)
	}
	f.PrintDefaults()
}

// PrintDefaults prints, to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set. See the
// documentation for the global function PrintDefaults for more information.
func (f *FlagSet) PrintDefaults() {
	f.VisitAll(func(flag *Flag) {
		fmt.Printf("  -%s %s\n", flag.Name, flag.Usage) // Two spaces before -; see next two comments.
	})
}
