package shared

import "github.com/go-cmd/cmd"

// Exec bundles stderr and stdout for a provided command and args, returning
// the status details and a possibly non-nil error
func Exec(logs chan string, c string, args ...string) (cmd.Status, error) {
	cmd := cmd.NewCmd(c, args...)
	outC := cmd.Stdout
	errC := cmd.Stderr

	go func() {
		for line := range outC {
			logs <- line
		}
	}()

	go func() {
		for line := range errC {
			logs <- line
		}
	}()

	status := <-cmd.Start()
	if status.Error != nil {
		return status, status.Error
	}

	return status, nil
}
