package exec

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"sync"
)

// LookPath is the wrapper of os/exec.LookPath
func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// RunCommandAndReturn runs a command, then returns the output
func RunCommandAndReturn(name, dir string, args ...string) (result string, err error) {
	stdout := &bytes.Buffer{}
	if err = RunCommandWithBuffer(name, dir, stdout, nil, args...); err == nil {
		result = stdout.String()
	}
	return
}

// RunCommandWithBuffer runs a command with buffer
// stdout and stderr could be nil
func RunCommandWithBuffer(name, dir string, stdout, stderr *bytes.Buffer, args ...string) error {
	if stdout == nil {
		stdout = &bytes.Buffer{}
	}
	if stderr != nil {
		stderr = &bytes.Buffer{}
	}
	return RunCommandWithIO(name, dir, stdout, stderr, args...)
}

// RunCommandWithIO runs a command with given IO
func RunCommandWithIO(name, dir string, stdout, stderr io.Writer, args ...string) (err error) {
	command := exec.Command(name, args...)
	if dir != "" {
		command.Dir = dir
	}

	//var stdout []byte
	//var errStdout error
	stdoutIn, _ := command.StdoutPipe()
	stderrIn, _ := command.StderrPipe()
	err = command.Start()
	if err != nil {
		return
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, _ = copyAndCapture(stdout, stdoutIn)
		wg.Done()
	}()

	_, _ = copyAndCapture(stderr, stderrIn)

	wg.Wait()

	err = command.Wait()
	return
}

// RunCommandInDir runs a command
func RunCommandInDir(name, dir string, args ...string) error {
	return RunCommandWithIO(name, dir, os.Stdout, os.Stderr, args...)
}

// RunCommand runs a command
func RunCommand(name string, arg ...string) (err error) {
	return RunCommandInDir(name, "", arg...)
}

// RunCommandWithSudo runs a command with sudo
func RunCommandWithSudo(name string, args ...string) (err error) {
	newArgs := make([]string, 0)
	newArgs = append(newArgs, name)
	newArgs = append(newArgs, args...)
	return RunCommand("sudo", newArgs...)
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
