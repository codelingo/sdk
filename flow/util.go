package flow

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	ctx := cli.NewContext(nil, fSet, nil)

	cmd.Action.(func(*cli.Context))(ctx)
	return nil
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
