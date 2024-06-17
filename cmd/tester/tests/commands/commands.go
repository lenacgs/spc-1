package commands

import (
	"fmt"
	"slices"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
)

type cmdbase struct {
	err error
}

type Command interface {
	Execute(cli.Client) error
	String() string
	GetError() error
}

func compareStatuses(expected []cli.Status, result []cli.Status) error {
	if !slices.Equal(expected, result) {
		return fmt.Errorf("result is invalid:\nResult: %v\nExpected: %v", result, expected)
	}

	return nil
}

func (b *cmdbase) GetError() error {
	return b.err
}
