package main

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
)

var (
	prompt        Prompt
	servers       *ServerList
	serverConfigs []*ServerConfig
	runCmd        string
)

func init() {
	flag.StringVar(&runCmd, "e", "", "Run command and exit")
	flag.Parse()

	prompt = NewPrompt("")

	serverConfigs = getServerConfigs()

	if runCmd == "" {
		fmt.Println("Connecting…")
	}

	servers = NewServerList(serverConfigs...)

	if runCmd == "" {
		fmt.Println("All servers connected")
	}
}

func main() {
	if runCmd != "" {
		if err := execInput(runCmd); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      prompt.String(),
		HistoryFile: path.Join(getHomeDir(), ".multishell"),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer rl.Close()

	mainLoop(rl)
}

func mainLoop(rl *readline.Instance) {
	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		if len(line) == 0 {
			continue
		}

		// Handle the execution of the input.
		if err = execInput(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if strings.HasPrefix(line, "cd") {
			servers.UpdatePwd(strings.TrimLeft(line, "cd "))
			prompt.SetPwd(servers.pwd)
			rl.SetPrompt(prompt.String())
		}
	}
}

func execInput(input string) error {
	input = strings.TrimSuffix(input, "\n")

	output := servers.Exec(input)
	os.Stdout.Write(output)

	return nil
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return "/tmp"
	}

	return usr.HomeDir
}

func getServerConfigs() []*ServerConfig {
	configFile := os.Getenv("SERVER_CONFIG")
	if configFile == "" {
		configFile = "/tmp/serverlist"
	}

	serverConfigFile, _ := ioutil.ReadFile(configFile)
	f := string(serverConfigFile)

	serverListRegex := regexp.MustCompile("SERVERS=\"(.*)\"")
	hostListRegex := regexp.MustCompile("HOSTS=\"(.*)\"")
	serverList := serverListRegex.FindStringSubmatch(f)
	hostList := hostListRegex.FindStringSubmatch(f)

	if len(serverList) > 1 && len(hostList) > 1 {
		serverList = strings.Split(serverList[1], " ")
		hostList = strings.Split(hostList[1], " ")
	}

	serverConfigs := make([]*ServerConfig, 0)
	for i, ip := range serverList {
		name := ip
		if len(hostList) >= i+1 {
			name = hostList[i]
		}

		serverConfigs = append(serverConfigs, &ServerConfig{
			Address: ip,
			Name:    name,
		})
	}

	return serverConfigs
}
