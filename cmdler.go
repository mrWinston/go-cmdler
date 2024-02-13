package cmdler

import (
	"bytes"
	"io"
	"os/exec"
)

type cmd struct {
	Stderr       string
	Stdout       string
	stdoutReader io.Reader
	stderrReader io.Reader
	Code         int
	Cmd          string
	prev         prevCmd
	config       *chainConfig
	Errors       []error
	stopChan     chan bool
	done         bool
}

type prevCmd interface {
	Run()
	HasRun() bool
	GetStdout() string
	GetStderr() string
	GetCode() int
	GetCommand() string
	GetErrors() []error
	HasErrors() bool
}

type chainConfig struct {
	// Stop execution of chain when one command has non-zero exit code
	StopOnErr bool
	// Shell to run command in.
	Shell string
	// Working dir of command. If empty, use pwd of parent
	Workdir string
}

func NewChainConfig() *chainConfig {
	return &chainConfig{
		StopOnErr: true,
		Shell:     "zsh",
		Workdir:   "",
	}
}

func New(command string) *cmd {
	return &cmd{
		Cmd:    command,
		config: NewChainConfig(),
	}
}

func (c *cmd) Pipe(command string) *cmd {
	return &cmd{
		Cmd:    command,
		config: NewChainConfig(),
		prev:   c,
	}
}

func (c *cmd) HasRun() bool {
	return c.done
}

func (c *cmd) Run() {
	if c.prev != nil {
		if !c.prev.HasRun() {
			c.prev.Run()
		}

		if c.config.StopOnErr && c.prev.GetCode() != 0 {
			c.Errors = append(c.Errors, NewPreviousCommandFailedError(c.prev))
			return
		}
	}

	exe := exec.Command(c.config.Shell, "-c", c.Cmd)
	exe.Dir = c.config.Workdir
	c.done = true

	stdin, err := exe.StdinPipe()
	if err != nil {
		c.Errors = append(c.Errors, err)
		c.Code = 128
		c.Stderr = err.Error()
		return
	}

	stdout, err := exe.StdoutPipe()
	if err != nil {
		c.Errors = append(c.Errors, err)
		c.Code = 128
		c.Stderr = err.Error()
		return
	}
	stderr, err := exe.StderrPipe()
	if err != nil {
		c.Errors = append(c.Errors, err)
		c.Code = 128
		c.Stderr = err.Error()
		return
	}

	go func() {
		defer stdin.Close()
		stdinString := ""
		if c.prev != nil {
			stdinString = c.prev.GetStdout()
		}
		_, err = io.WriteString(stdin, stdinString)
		if err != nil {
			c.Errors = append(c.Errors, err)
		}
	}()

	err = exe.Start()
	if err != nil {
		c.Errors = append(c.Errors, err)
		c.Code = 128
		c.Stderr = err.Error()
	}

	c.Stdout = ReaderToString(stdout)
	c.Stderr = ReaderToString(stderr)

	if err := exe.Wait(); err != nil {
		eerr, ok := err.(*exec.ExitError)
		if ok {
			c.Code = eerr.ExitCode()
		} else {
			c.Code = 128
			c.Stderr = err.Error()
		}
	} else {
		c.Code = 0
	}
}

func (c *cmd) Redirect(writer io.Writer) {
}

func ReaderToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

func (c *cmd) HasErrors() bool {
	return len(c.Errors) != 0
}

func (c *cmd) GetStdout() string {
	return c.Stdout
}
func (c *cmd) GetStderr() string {
	return c.Stderr
}
func (c *cmd) GetCode() int {
	return c.Code
}
func (c *cmd) GetCommand() string {
	return c.Cmd
}
func (c *cmd) GetErrors() []error {
	return c.Errors
}
