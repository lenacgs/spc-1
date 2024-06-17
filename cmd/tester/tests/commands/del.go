package commands

import (
	"fmt"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
)

type del struct {
	cmdbase
	ids      []cli.Id
	expected []cli.Status
}

func Del(ids []cli.Id, old []cli.Status) Command {
	var err error

	if len(ids) != len(old) {
		err = fmt.Errorf("del's old status size different than ids size")
	}

	return &del{
		cmdbase: cmdbase{
			err: err,
		},
		ids:      ids,
		expected: old,
	}
}

func (d *del) Execute(c cli.Client) error {
	if d.err != nil {
		return d.err
	}

	result, err := c.Del(d.ids)
	if err != nil {
		d.err = err
		return err
	}

	d.err = compareStatuses(d.expected, result)
	return d.err
}

func (d *del) String() string {
	return fmt.Sprintf("del %v %v", d.ids, d.expected)
}
