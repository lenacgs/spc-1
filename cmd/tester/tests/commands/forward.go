package commands

import (
	"fmt"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
)

type forward struct {
	cmdbase
	t cli.Time
}

func Forward(t cli.Time) Command {
	var err error

	if t < 0 {
		err = fmt.Errorf("negative time is not allowed")
	}

	return &forward{
		cmdbase: cmdbase{
			err: err,
		},
		t: t,
	}
}

func (f *forward) Execute(c cli.Client) error {
	if f.err != nil {
		return f.err
	}

	f.err = c.Forward(f.t)
	return f.err
}

func (f *forward) String() string {
	return fmt.Sprintf("forward %v", f.t)
}
