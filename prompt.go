package main

import (
	"fmt"
	"github.com/fatih/color"
)

type (
	Prompt string
)

func NewPrompt(pwd string) Prompt {
	userString := color.New(color.FgCyan)
	serverString := color.New(color.FgYellow)
	pathString := color.New(color.Bold)

	return Prompt(fmt.Sprintf(
		"%s @ %s in %s > ",
		userString.Sprint("root"),
		serverString.Sprint("servers"),
		pathString.Sprint(pwd),
	))
}

func (p Prompt) String() string {
	return string(p)
}

func (p *Prompt) SetPwd(pwd string) {
	*p = NewPrompt(pwd)
}
