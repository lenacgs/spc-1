package tests

import (
	"strings"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
	cmd "gitlab.com/sibsfps/spc/spc-1/cmd/tester/tests/commands"
)

type testbase struct {
	batch []cmd.Command
}

type TestFunc func(cli.Client) Test

type Test interface {
	Execute(client cli.Client) error
	String(showErrors ...bool) string
}

func (t *testbase) Execute(client cli.Client) error {
	var err error
	for _, c := range t.batch {
		err = c.Execute(client)
		if err != nil {
			break
		}
	}

	return err
}

func (t *testbase) String(showErrors ...bool) string {
	var sb strings.Builder
	for _, c := range t.batch {
		sb.WriteString(c.String())
		sb.WriteString("\n")
		if len(showErrors) > 0 && showErrors[0] {
			err := c.GetError()
			if err != nil {
				sb.WriteString("^^^^^^^^^^^^^^^^\n")
				sb.WriteString(err.Error())
				sb.WriteString("\n================\n")
			}
		}
	}

	return sb.String()
}
