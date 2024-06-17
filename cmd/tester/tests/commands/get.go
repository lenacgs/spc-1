package commands

import (
	"fmt"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
)

type get struct {
	cmdbase
	ids      []cli.Id
	expected []cli.Status
}

func Get(ids []cli.Id, cur []cli.Status) Command {
	var err error

	if len(ids) != len(cur) {
		err = fmt.Errorf("del's old status size different than ids size")
	}

	return &get{
		cmdbase: cmdbase{
			err: err,
		},
		ids:      ids,
		expected: cur,
	}
}

func (g *get) Execute(c cli.Client) error {
	if g.err != nil {
		return g.err
	}

	result, err := c.Get(g.ids)
	if err != nil {
		g.err = err
		return err
	}

	g.err = compareStatuses(g.expected, result)
	return g.err
}

func (g *get) String() string {
	return fmt.Sprintf("get %v %v", g.ids, g.expected)
}
