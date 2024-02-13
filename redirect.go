package cmdler

import (
	"fmt"
	"io"
)

type redirect struct {
	prev   *cmd
	target io.Writer
	done   bool
}

func (r *redirect) Run() {
	if r.prev != nil {
		if !r.prev.HasRun() {
			r.prev.Run()
		}
	}
	fmt.Fprint(r.target, r.prev.Stdout)
}

type inputRedirect struct {
	stdout string
}

func NewStaticInput(in io.Reader) *inputRedirect {
	return &inputRedirect{
		stdout: ReaderToString(in),
	}
}

// GetCode implements prevCmd.
func (i *inputRedirect) GetCode() int {
	return 0
}

// GetCommand implements prevCmd.
func (i *inputRedirect) GetCommand() string {
	return ""
}

// GetErrors implements prevCmd.
func (i *inputRedirect) GetErrors() []error {
	return nil
}

// GetStderr implements prevCmd.
func (i *inputRedirect) GetStderr() string {
	return ""
}

// GetStdout implements prevCmd.
func (i *inputRedirect) GetStdout() string {
	return i.stdout
}

// HasErrors implements prevCmd.
func (*inputRedirect) HasErrors() bool {
	return false
}

// HasRun implements prevCmd.
func (*inputRedirect) HasRun() bool {
	return true
}

// Run implements prevCmd.
func (*inputRedirect) Run() {}
