package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"

	lamport "github.com/ISSuh/logical-clock"
	"github.com/chzyer/readline"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("connect"),
	readline.PcItem("get"),
	readline.PcItem("put"),
	readline.PcItem("del"),
	readline.PcItem("help"),
	readline.PcItem("service"),
	readline.PcItem("cache"))

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

func split(s string, sep rune) []string {
	return strings.FieldsFunc(s,
		func(c rune) bool {
			return c == sep
		},
	)
}

func parseInts(tokens []string, skip int) ([]int, error) {
	var ints []int = make([]int, 0)
	l := len(tokens)
	for i := 0; i < l; i += 1 + skip {
		intToken, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse WorkerID %s", tokens[i])
		}

		ints = append(ints, int(intToken))
	}

	return ints, nil
}

func checkConnected(client *RestClient) bool {
	if client == nil {
		log.Println("Not connected")
		return false
	}
	return true
}

func connect(cmd []string) *RestClient {
	if len(cmd) == 1 {
		cmd = append(cmd, "http://localhost:8080")
	}

	if len(cmd) != 2 {
		log.Println("Use: connect {url}")
		log.Println("Example: connect http://localhost:8080")
		return nil
	}
	urlString := cmd[1]

	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = fmt.Sprintf("http://%s", urlString)
	}

	clientUrl, err := url.Parse(urlString)
	if err != nil {
		log.Println("Warning couldn't parse URL", strconv.Quote(err.Error()))
	}

	restClient, err := MakeRestClient(*clientUrl)
	if err != nil {
		log.Println("Warning", err.Error())
		return restClient
	}
	log.Printf("Connected to %s\n", restClient.serverURL.String())

	return restClient
}

func get(client *RestClient, cmd []string) {
	if !checkConnected(client) {
		return
	}

	args := cmd[1:]

	workers, err := parseInts(args, 0)
	if err != nil {
		log.Println(err)
	}

	res, err := client.Get(workers)
	if err != nil {
		log.Println("Error", err.Error())
		return
	}

	for n, worker := range workers {
		log.Printf("Worker %d is %d\n", worker, res[n].New)
	}
}

func put(client *RestClient, cmd []string) {
	if !checkConnected(client) {
		return
	}

	args := cmd[1:]

	if (len(args) % 2) != 0 {
		log.Printf("Invalid number of args")
		log.Printf("Use: put [{worker} {status}]...")
		return
	}

	workers, err := parseInts(args, 1)
	if err != nil {
		log.Println(err)
	}
	statuses, err := parseInts(args[1:], 1)
	if err != nil {
		log.Println(err)
	}

	res, err := client.Put(workers, statuses)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for n, worker := range workers {
		log.Printf("Worker %d is now %d and was %d\n", worker, res[n].New, res[n].Old)
	}
}

func del(client *RestClient, cmd []string) {
	if !checkConnected(client) {
		return
	}

	args := cmd[1:]

	workers, err := parseInts(args, 0)
	if err != nil {
		log.Println(err)
	}

	res, err := client.Del(workers)
	if err != nil {
		log.Println("Error", err.Error())
		return
	}

	for n, worker := range workers {
		log.Printf("Worker %d was %d\n", worker, res[n].Old)
	}
}

func cache(client *RestClient, cmd []string, clock *lamport.LamportClock) {
	if !checkConnected(client) {
		return
	}

	args := cmd[1:]

	workers, err := parseInts(args, 0)
	if err != nil {
		log.Println("Error", err)
	}

	res, err := client.Cache(workers, clock)

	if err != nil {
		log.Println("Error", err.Error())
		return
	}

	for n, worker := range workers {
		log.Printf("Worker %d is %d\n", worker, res[n].Status)
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "client.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	var restClient *RestClient
	var clock = lamport.NewLamportClock()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		cmd := split(line, ' ')

		switch {
		case len(cmd) == 0:
		case cmd[0] == "connect":
			restClient = connect(cmd)

		case cmd[0] == "get":
			get(restClient, cmd)

		case cmd[0] == "put":
			put(restClient, cmd)

		case cmd[0] == "del":
			del(restClient, cmd)

		case cmd[0] == "help":
			usage(l.Stderr())

		case cmd[0] == "cache":
			cache(restClient, cmd, clock)

		default:
			log.Println("Unknown command", strconv.Quote(line))

		}
	}
}
