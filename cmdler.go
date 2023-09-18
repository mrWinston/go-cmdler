package main

import (
	"bytes"
	"io"
	"os/exec"
)

type cmd struct {
	Stderr string
	Stdout string
	Code   int
	Cmd    string
	prev   *cmd
	config *chainConfig
	Errors []error
}

type chainConfig struct {
	// Stop execution of chain when one command has non-zero exit code
	StopOnErr bool
	// Wether or not to redirect Stderr to Stdout
	RedirStderr bool
	// Shell to run command in.
	Shell string
	// Working dir of command. If empty, use pwd of parent
	Workdir string
}

func NewChainConfig() *chainConfig {
	return &chainConfig{
		StopOnErr:   true,
		RedirStderr: false,
		Shell:       "zsh",
		Workdir:     "",
	}
}

func NewChain(conf *chainConfig) *cmd {
	return &cmd{
		Code:   0,
		Cmd:    "",
		prev:   nil,
		config: conf,
		Errors: []error{},
	}
}

func ReaderToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

func (c *cmd) HasErrors() bool {
	return len(c.Errors) != 0
}

func (c *cmd) New(command string) *cmd {
	out := &cmd{}
	out.config = c.config
	out.prev = c

	// quick exit
	if c.config.StopOnErr && c.Code != 0 {
		out.Stdout = c.Stdout
		out.Stderr = c.Stderr
		out.Code = c.Code
		return out
	}

	exe := exec.Command(c.config.Shell, "-c", command)
	exe.Dir = c.config.Workdir

	stdin, err := exe.StdinPipe()
	if err != nil {
		out.Errors = append(out.Errors, err)
		out.Code = 128
		out.Stderr = err.Error()
		return out
	}

	stdout, err := exe.StdoutPipe()
	if err != nil {
		out.Errors = append(out.Errors, err)
		out.Code = 128
		out.Stderr = err.Error()
		return out
	}
	stderr, err := exe.StderrPipe()
	if err != nil {
		out.Errors = append(out.Errors, err)
		out.Code = 128
		out.Stderr = err.Error()
		return out
	}

	go func() {
		defer stdin.Close()
		_, err = io.WriteString(stdin, c.Stdout)
		if err != nil {
			out.Errors = append(out.Errors, err)
		}
	}()

	err = exe.Start()
	if err != nil {
		out.Errors = append(out.Errors, err)
		out.Code = 128
		out.Stderr = err.Error()
		return out
	}

	out.Stdout = ReaderToString(stdout)
	out.Stderr = ReaderToString(stderr)

	if err := exe.Wait(); err != nil {
		eerr, ok := err.(*exec.ExitError)
		if ok {
			out.Code = eerr.ExitCode()
		} else {
			out.Code = 128
			out.Stderr = err.Error()
		}
	} else {
		out.Code = 0
	}
	return out
}
