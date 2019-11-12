package main

import (
	"bufio"
	"context"
	"os/exec"

	"github.com/pieterclaerhout/go-log"
	"golang.org/x/sync/errgroup"
)

func main() {

	// Print the log timestamps
	log.PrintTimestamp = true

	// The command you want to run along with the argument
	cmd := exec.Command("brew", "info", "golang")

	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Start the command and check for errors
	err := cmd.Start()
	log.CheckError(err)

	// Start a new error group, run the command and wait for it to finish
	// We are running this in a go routine so that we don't block
	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		return cmd.Wait()
	})

	// Use the scanner to scan the output line by line and log ti
	for scanner.Scan() {
		line := scanner.Text()
		log.Info(line)
	}

	// Wait for the command to finish and check for errors
	err = errs.Wait()
	log.CheckError(err)

}
