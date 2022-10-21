package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/helloyi/go-sshclient"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os/user"
	"path"
	"time"
)

type (
	Server struct {
		config    *ServerConfig
		client    *sshclient.Client
		Connected bool
		pwd       string
	}
	ServerConfig struct {
		Address string
		Name    string
	}
)

func NewServer(config *ServerConfig) (*Server, error) {
	sshKey, err := sshKeyPath()
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(sshKey)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		Timeout:         10 * time.Second,
	}

	client, err := sshclient.Dial("tcp", fmt.Sprintf("%s:22", config.Address), sshConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		client:    client,
		Connected: false,
		config:    config,
	}, nil
}

func (s *Server) Exec(cmd string) ([]byte, error) {
	out, err := s.client.Cmd("cd " + s.pwd + " && " + cmd).Output()
	if err != nil {
		errFormat := color.New(color.Bold).Add(color.FgRed)
		prefix := []byte(errFormat.Sprintf("[%s] ", s.config.Name))
		return nil, fmt.Errorf("%s %s\n", string(prefix), err.Error())
	}

	bold := color.New(color.Bold).Add(color.FgHiBlack)
	prefix := []byte(bold.Sprintf("[%s] ", s.config.Name))
	tmp := bytes.Split(out, []byte("\n"))
	for i, a := range tmp {
		if len(a) != 0 {
			tmp[i] = append(prefix, a...)
		}
	}

	return bytes.Join(tmp, []byte("\n")), nil
}

func sshKeyPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(usr.HomeDir, ".ssh", "id_rsa"), nil
}
