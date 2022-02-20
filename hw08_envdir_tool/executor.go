package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	cmdRun := exec.Command(cmd[0], cmd[1:]...)
	var cmdEnv []string
	for en := range env {
		cmdEnv = append(cmdEnv, en+"="+env[en].Value)
	}
	cmdRun.Env = cmdEnv
	cmdRun.Stdin = os.Stdin
	cmdRun.Stdout = os.Stdout

	if err := cmdRun.Run(); err != nil {
		log.Fatal(err)
	}
	return cmdRun.ProcessState.ExitCode()
}
