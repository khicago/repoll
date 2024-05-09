package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const (
	prefix = "\033[94m||\033[0m"
)

func controlPrintf(format string, a ...any) {
	fmt.Printf("\033[94m"+format+"\033[0m", a...)
}

func printPrefixedOutput(pipe io.ReadCloser, prefix string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		// \033[90m sets the color to bright black (grey), and \033[0m resets the color
		fmt.Printf("\t%s \033[90m%s\033[0m\n", prefix, scanner.Text())
	}
}

func runCommand(cmd *exec.Cmd) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	go printPrefixedOutput(stdoutPipe, prefix)
	go printPrefixedOutput(stderrPipe, prefix)

	err = cmd.Wait()
	return err
}

func runCommandWithTimer(cmd *exec.Cmd) error {
	cmdDone := make(chan struct{})
	startTime := time.Now()

	// Start the spinner in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		spinnerPosition := 0
		for {
			since := int(time.Since(startTime).Milliseconds())
			strSince := strconv.Itoa(since) + "ms"
			if since >= 1000 {
				strSince = strconv.Itoa(since/1000) + "." + strconv.Itoa((since%1000)/100) + "s"
			}
			select {
			case <-cmdDone:
				controlPrintf("\r[%s] %s Command finished.\n", strSince, prefix)
				return
			default:
				controlPrintf("\r[%s]", strSince)
				time.Sleep(100 * time.Millisecond)
				spinnerPosition++
			}
		}
	}()

	controlPrintf("[0ms]\t%s Command `%s` Started.\n", prefix, cmd.String())
	// Run the command and wait for it to finish
	err := runCommand(cmd)

	// Signal spinner to stop and wait for it
	close(cmdDone)
	wg.Wait()

	return err
}
