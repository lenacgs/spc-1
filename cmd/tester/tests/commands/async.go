package commands

import (
	"fmt"
	"strings"
	"sync"

	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
)

type async struct {
	cmdbase
	id    string
	cmd   []Command
	wg    sync.WaitGroup
	errAt int
}

type await struct {
	cmdbase
	ids   []string
	errAt int
}

var asyncIds map[string]*async = make(map[string]*async)

func Async(id string, cmd ...Command) Command {
	var err error = nil
	_, ok := asyncIds[id]
	if ok {
		err = fmt.Errorf("async id already in use")
	}

	a := &async{
		id:  id,
		cmd: cmd,
		wg:  sync.WaitGroup{},
		cmdbase: cmdbase{
			err: err,
		},
	}

	asyncIds[id] = a

	return a
}

func (a *async) Execute(c cli.Client) error {
	if a.err != nil {
		return a.err
	}

	a.wg.Add(1)
	go func(c cli.Client, a *async) {
		defer a.wg.Done()

		for i := 0; i < len(a.cmd); i++ {
			a.err = a.cmd[i].Execute(c)
			if a.err != nil {
				a.errAt = i
				break
			}
		}
	}(c, a)

	return nil
}

func (a *async) String() string {
	sb := strings.Builder{}
	sb.WriteString(
		fmt.Sprintf("async \"%s\"", a.id),
	)

	if len(a.cmd) == 0 {
		return sb.String()
	}
	if len(a.cmd) == 1 {
		return fmt.Sprintf("%s %s", sb.String(), a.cmd[0].String())
	}

	sb.WriteString(":")
	for _, cmd := range a.cmd {
		sb.WriteString(
			fmt.Sprintf("\n  + %s", strings.ReplaceAll(cmd.String(), "\n", "\n  ")),
		)
	}

	return sb.String()
}

func Await(ids ...string) Command {
	aw := &await{
		ids: ids,
	}

	for i := 0; i < len(ids); i++ {
		as, ok := asyncIds[ids[i]]
		if ok {
			aw.err = as.err
		} else {
			aw.err = fmt.Errorf("async id does not exist")
		}

		if aw.err != nil {
			aw.errAt = i
			break
		}
	}

	return aw
}

func (a *await) Execute(c cli.Client) error {
	for i := 0; i < len(a.ids); i++ {
		id := a.ids[i]
		as, ok := asyncIds[id]
		if !ok {
			a.errAt = i
			a.err = fmt.Errorf("async id does not exist")
			return a.err
		}
		as.wg.Wait()
		delete(asyncIds, id)

		if as.err != nil {
			a.errAt = i
			a.err = fmt.Errorf("check async \"%s\"", a.ids[i])
			break
		}
	}

	return a.err
}

func (a *await) String() string {
	sb := strings.Builder{}
	sb.WriteString("await")

	if len(a.ids) == 0 {
		return sb.String()
	}
	if len(a.ids) == 1 {
		return fmt.Sprintf("%s \"%s\"", sb.String(), a.ids[0])
	}

	sb.WriteString(":")
	for _, id := range a.ids {
		sb.WriteString(
			fmt.Sprintf("\n  - \"%s\"", id),
		)
	}

	return sb.String()
}
