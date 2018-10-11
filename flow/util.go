package flow

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

// TODO(waigani) move this to codelingo/sdk/flow
func Run(cmd cli.Command) error {
	fSet := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	for _, flag := range cmd.Flags {
		flag.Apply(fSet)
	}

	if err := fSet.Parse(os.Args[1:]); err != nil {
		return errors.Trace(err)
	}
	if err := normalizeFlags(cmd.Flags, fSet); err != nil {
		return errors.Trace(err)
	}

	ctx := cli.NewContext(nil, fSet, nil)

	cmd.Action.(func(*cli.Context))(ctx)
	return nil
}

func normalizeFlags(flags []cli.Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		parts := strings.Split(f.GetName(), ",")
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if !visited[name] {
				copyFlag(name, ff, set)
			}
		}
	}
	return nil
}

func copyFlag(name string, ff *flag.Flag, set *flag.FlagSet) {
	switch ff.Value.(type) {
	case *cli.StringSlice:
	default:
		set.Set(name, ff.Value.String())
	}
}

// TODO(waigani) move this to codelingo/sdk/flow
func HandleErr(err error) {
	if errors.Cause(err).Error() == "ui" {
		if e, ok := err.(*errors.Err); ok {
			log.Println(e.Underlying())
			fmt.Println(e.Underlying())
			os.Exit(1)
		}
	}
	fmt.Println(err.Error())
}

// TODO(waigani) this should live under the VCS domain, not Flows
const NoCommitErrMsg = "This looks like a new repository. Please make an initial commit before running `lingo review`. This is only Required for the initial commit, subsequent changes to your repo will be picked up by lingo without committing."

// TODO(waigani) use typed error
func NoCommitErr(err error) bool {
	return strings.Contains(err.Error(), "ambiguous argument 'HEAD'")
}
