package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

type (
	ServerList struct {
		list []*Server
		pwd  string
	}
)

func NewServerList(servers ...*ServerConfig) *ServerList {
	sl := make([]*Server, len(servers))

	wg := sync.WaitGroup{}
	wg.Add(len(servers))

	for i, config := range servers {
		go func(i int, config *ServerConfig) {
			server, err := NewServer(config)

			if err != nil {
				wg.Done()
				os.Stdout.Write([]byte(fmt.Sprintf("Could not connect to %s\n", config.Name)))
				return
			}

			wg.Done()
			sl[i] = server
		}(i, config)
	}

	wg.Wait()

	serverlist := &ServerList{
		list: sl,
	}

	serverlist.UpdatePwd("/root")
	prompt.SetPwd(serverlist.pwd)

	return serverlist
}

func (sl *ServerList) Exec(cmd string) []byte {
	wg := sync.WaitGroup{}
	wg.Add(len(sl.list))

	outputs := make([][]byte, len(sl.list))

	for i, server := range sl.list {
		go func(i int, cmd string, server *Server) {
			output, err := server.Exec(cmd)
			if err != nil {
				outputs[i] = []byte(err.Error())
			} else {
				outputs[i] = output
			}
			wg.Done()
		}(i, cmd, server)
	}

	wg.Wait()

	filteredOutputs := make([][]byte, 0, len(outputs))
	for _, output := range outputs {
		if len(output) != 0 {
			filteredOutputs = append(filteredOutputs, output)
		}
	}

	return bytes.Join(filteredOutputs, []byte("\n"))
}

func (sl *ServerList) UpdatePwd(pwd string) {
	(*sl).pwd = pwd
	for _, server := range (*sl).list {
		server.pwd = pwd
	}
}
