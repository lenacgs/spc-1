package tests

import (
	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
	cmd "gitlab.com/sibsfps/spc/spc-1/cmd/tester/tests/commands"
)

func Test1(client cli.Client) Test {
	return &testbase{
		batch: []cmd.Command{
			cmd.Delay(0),

			cmd.Get([]cli.Id{0, 1}, []cli.Status{0, 0}),

			cmd.Put(
				[]cli.Record{
					{Id: 2, Status: 1},
					{Id: 3, Status: 2},
				},
			),
			cmd.Get([]cli.Id{2, 3}, []cli.Status{1, 2}),

			cmd.Forward(client.SoftTTL()),
			cmd.Put(
				[]cli.Record{
					{Id: 0, Status: 1},
					{Id: 1, Status: 2},
				},
			),
			cmd.Get([]cli.Id{0, 1, 2, 3}, []cli.Status{1, 2, 1, 2}),

			cmd.Forward(client.HardTTL() - client.SoftTTL()),
			cmd.Del([]cli.Id{2, 3}, []cli.Status{1, 2}),
			cmd.Get([]cli.Id{0, 1, 2, 3}, []cli.Status{1, 2, 0, 0}),

			cmd.Forward(client.SoftTTL()),
			cmd.Del([]cli.Id{0, 1}, []cli.Status{1, 2}),
			cmd.Put(
				[]cli.Record{
					{Id: 2, Status: 2},
					{Id: 3, Status: 1},
				},
			),
			cmd.Get([]cli.Id{0, 1, 2, 3}, []cli.Status{0, 0, 2, 1}),
		},
	}
}

/*
delay 0

get [0,1] [0,0]

put [(2,1),(3,2)]
get [2,3] [1,2]

forward soft-ttl
put [(0,1), (0,2)]
get [0,1,2,3] [1,2,1,2]

forward - hard-ttl soft-ttl
get [0,1,2,3] [1,2,0,0]

forward soft-ttl
put [(2,2),(3,1)]
get [0,1,2,3] [0,0,2,1]
*/
