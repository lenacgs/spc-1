package commands

import (
	"fmt"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
)

type delay struct {
	cmdbase
	t cli.Time
}

func Delay(t cli.Time) Command {
	var err error

	if t < 0 {
		err = fmt.Errorf("negative time is not allowed")
	}

	return &delay{
		cmdbase: cmdbase{
			err: err,
		},
		t: t,
	}
}

func (d *delay) Execute(c cli.Client) error {
	if d.err != nil {
		return d.err
	}

	d.err = c.Delay(d.t)
	return d.err
}

func (d *delay) String() string {
	return fmt.Sprintf("delay %v", d.t)
}
