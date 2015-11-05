package tryprog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Tryprog tries the program at "path" and compares its output to "expected")
//func Try(expected []byte, path string, args []string) (err error) {
func Try(expected []byte, path string) (err error) {

	//cmd := exec.Command(path, args...)
	cmd := exec.Command(path)
	defer os.Remove(path)
	var out []byte

	c1 := make(chan string, 1)
	go func() {
		out, err = cmd.CombinedOutput()
		c1 <- "success" // dummy
	}()

	select {
	case <-time.After(time.Second * 1):
		{
			if err := cmd.Process.Kill(); err != nil {
				panic(fmt.Sprintf("failed to kill: %s", err))
			}
			err = fmt.Errorf("time-out-no-exit")
		}

	case <-c1:
		switch {
		case (bytes.Compare(expected, out) != 0) && (err != nil):
			{
				err = fmt.Errorf("bad-out-bad-exit")
			}
		case bytes.Compare(expected, out) != 0:
			{
				err = fmt.Errorf("bad-out-good-exit")
			}
		case err != nil:
			{
				err = fmt.Errorf("good-out-bad-exit")
			}
		}

	}
	return
}
