package commands

import (
	"fmt"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
)

type put struct {
	cmdbase
	records  []cli.Record
	expected []cli.Status
}

func Put(records []cli.Record) Command {
	var expected []cli.Status
	for _, r := range records {
		expected = append(expected, cli.Status(r.Status))
	}

	return &put{
		cmdbase: cmdbase{
			err: nil,
		},
		records:  records,
		expected: expected,
	}
}

func (p *put) Execute(c cli.Client) error {
	if p.err != nil {
		return p.err
	}

	result, err := c.Put(p.records)
	if err != nil {
		p.err = err
		return err
	}

	p.err = compareStatuses(p.expected, result)
	return p.err
}

func (p *put) String() string {
	return fmt.Sprintf("put %v", p.records)
}
