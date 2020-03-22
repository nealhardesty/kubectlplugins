package global

import (
	"flag"
)

type GlobalOptions struct {
	UseGlobalRegex bool
	NonInteractive bool
	PrintOnly      bool
}

func NewGlobalOptionsWithFlag() *GlobalOptions {
	options := GlobalOptions{}

	flag.BoolVar(&options.NonInteractive, "d", false, "force non interactive mode, pick sensible defaults (shorthand)")
	flag.BoolVar(&options.NonInteractive, "-non-interactive", false, "force non interactive mode, pick sensible defaults")
	flag.BoolVar(&options.UseGlobalRegex, "g", false, "regex searches can be found anywhere in the match, not just at the start (shorthand)")
	flag.BoolVar(&options.UseGlobalRegex, "-global-regex", false, "regex searches can be found anywhere in the match, not just at the start")
	flag.BoolVar(&options.PrintOnly, "p", false, "print the possible values and exit (shorthand)")
	flag.BoolVar(&options.PrintOnly, "-print", false, "print the possible values and exit (shorthand)")

	return &options
}
