package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/sibsfps/spc/spc-1/cmd/tester/clients"
	"gitlab.com/sibsfps/spc/spc-1/cmd/tester/tests"
)

var serviceHost = flag.String("s", "", "ip:port of Service")
var workersHost = flag.String("w", "", "ip:port of Workers")

func main() {
	flag.Parse()

	s, err := clients.NewService(*serviceHost)
	if err != nil {
		fmt.Printf("Could not create service client: %v", err)
		os.Exit(1)
	}

	w, err := clients.NewWorkers(*workersHost)
	if err != nil {
		fmt.Printf("Could not create worker client: %v", err)
		os.Exit(1)
	}

	c := clients.NewClient(
		100,
		1000,
		s,
		w,
	)

	batch := []tests.TestFunc{
		tests.Test1,
		tests.Test2,
	}

	for i, tf := range batch {
		t := tf(c)
		err := t.Execute(c)
		res := "Ok"
		if err != nil {
			res = "Failed"
		}
		fmt.Printf("Test[%d]: %s\n", i, res)
		if err != nil {
			fmt.Println(t.String(true))
		}
		c.Forward(c.HardTTL())
	}
}
