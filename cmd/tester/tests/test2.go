package tests

import (
	cli "gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
	cmd "gitlab.com/sibsfps/spc/spc-1/cmd/tester/tests/commands"
)

func Test2(client cli.Client) Test {
	return &testbase{
		batch: []cmd.Command{
			cmd.Delay(0),

			cmd.Async(
				"1",
				cmd.Get([]cli.Id{0, 1}, []cli.Status{0, 0}),
			),

			cmd.Async(
				"2",
				cmd.Put(
					[]cli.Record{
						{Id: 2, Status: 1},
						{Id: 3, Status: 2},
					},
				),
				cmd.Get([]cli.Id{2, 3}, []cli.Status{1, 2}),
			),

			cmd.Async(
				"3",
				cmd.Await("1", "2"),
				cmd.Async(
					"3.1",
					cmd.Forward(client.SoftTTL()),
				),
				cmd.Async(
					"3.2",
					cmd.Put(
						[]cli.Record{
							{Id: 0, Status: 1},
							{Id: 1, Status: 2},
						},
					),
				),
				cmd.Await("3.1", "3.2"),
				cmd.Get([]cli.Id{0, 1, 2, 3}, []cli.Status{1, 2, 1, 2}),
			),

			cmd.Async(
				"4",
				cmd.Await("3"),
				cmd.Async(
					"4.1",
					cmd.Del([]cli.Id{2, 3}, []cli.Status{1, 2}),
				),
				cmd.Async(
					"4.2",
					cmd.Forward(client.HardTTL()-client.SoftTTL()),
				),
				cmd.Await("4.1", "4.2"),
			),

			cmd.Async(
				"5",
				cmd.Await("4"),
				cmd.Get([]cli.Id{0, 1, 2, 3}, []cli.Status{1, 2, 0, 0}),
				cmd.Del([]cli.Id{0, 1}, []cli.Status{1, 2}),
			),

			cmd.Async(
				"6",
				cmd.Await("5"),
				cmd.Async(
					"6.1",
					cmd.Forward(client.SoftTTL()),
				),
				cmd.Async(
					"6.2",
					cmd.Put(
						[]cli.Record{
							{Id: 2, Status: 2},
							{Id: 3, Status: 1},
						},
					),
				),
				cmd.Await("6.1", "6.2"),
				cmd.Get([]cli.Id{0, 1, 2, 3}, []cli.Status{0, 0, 2, 1}),
			),

			cmd.Await("6"),
		},
	}
}
