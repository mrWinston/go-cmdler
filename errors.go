package cmdler

import "fmt"

type PreviousCommandFailedError struct {
	PreviousCommand string
	Code            int
	Stderr          string
}

func (p *PreviousCommandFailedError) Error() string {
	return fmt.Sprintf(
		"Previous Command '%s' failed with Error Code '%d'. Error Output: %s",
		p.PreviousCommand,
		p.Code,
		p.Stderr,
	)
}

func NewPreviousCommandFailedError(c prevCmd) error {
	return &PreviousCommandFailedError{
		PreviousCommand: c.GetCommand(),
		Code:            c.GetCode(),
		Stderr:          c.GetStderr(),
	}
}
