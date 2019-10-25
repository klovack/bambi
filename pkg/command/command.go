package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// IsAvailable returns true if the given name(command) exists in user's computer
func IsAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Execute executes command and output anything to the console
// Returns error from the cmd.Run()
func Execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to run command %s because: %v ", name, err)
	}

	return nil
}

// Executef executes command and return the output as string
// Returns string
func Executef(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to run %s because: %v ", name, err)
	}

	return out.String(), nil
}

// ExecutePipe takes the output buffer and slice of exec.Cmd and executes it with pipe.
// It returns error if the pipe is broken or there's an error while executing
func ExecutePipe(outputBuffer *bytes.Buffer, stack ...*exec.Cmd) (err error) {
	var errorBuffer bytes.Buffer
	pipeStack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdinPipe, stdoutPipe := io.Pipe()
		stack[i].Stdout = stdoutPipe
		stack[i].Stderr = &errorBuffer
		stack[i+1].Stdin = stdinPipe
		pipeStack[i] = stdoutPipe
	}
	stack[i].Stdout = outputBuffer
	stack[i].Stderr = &errorBuffer

	if err := callPipe(stack, pipeStack); err != nil {
		return fmt.Errorf(string(errorBuffer.Bytes()), err)
	}
	return err
}

func callPipe(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = callPipe(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}
