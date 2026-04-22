package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

func parsePorts(s string) []string {
	var ports []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			ports = append(ports, p)
		}
	}
	return ports
}

func main() {
	local := flag.String("local", "", "Local port forwarding: comma-separated list of ports")
	remote := flag.String("remote", "", "Remote port forwarding: comma-separated list of ports")
	flag.Parse()

	if *local == "" && *remote == "" {
		fmt.Fprintln(os.Stderr, "Error: at least one of -local or -remote must be specified")
		os.Exit(1)
	}

	sshArgs := flag.Args()
	if len(sshArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Error: ssh arguments must be specified after --")
		os.Exit(1)
	}

	var args []string

	if *local != "" {
		for _, port := range parsePorts(*local) {
			args = append(args, "-L", fmt.Sprintf("%s:127.0.0.1:%s", port, port))
		}
	}

	if *remote != "" {
		for _, port := range parsePorts(*remote) {
			args = append(args, "-R", fmt.Sprintf("%s:127.0.0.1:%s", port, port))
		}
	}

	args = append(args, sshArgs...)

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Fprintf(os.Stderr, "Executing: ssh %s\n", strings.Join(args, " "))

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- cmd.Wait()
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			fmt.Fprintf(os.Stderr, "SSH connection closed: %v\n", err)
			os.Exit(1)
		}
	case <-sigCh:
		fmt.Fprintln(os.Stderr, "\nKeyboard interrupt received. Terminating SSH connection.")
		_ = cmd.Process.Kill()
		select {
		case <-doneCh:
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
		}
		fmt.Fprintln(os.Stderr, "SSH connection closed. Exiting wrapper.")
		os.Exit(0)
	}
}
