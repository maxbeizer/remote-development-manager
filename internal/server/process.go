package server

import (
	"os/exec"
	"sync"
)

type processManager struct {
	commands map[int]*exec.Cmd
	mu       sync.Mutex
}

func (pm *processManager) Start(name string, path string, args ...string) error {
	cmd := exec.Command(path, args...)
	err := cmd.Start()

	if err != nil {
		return err
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()
	pid := cmd.Process.Pid
	pm.commands[pid] = cmd

	go func(cmd *exec.Cmd, pid int) {
		cmd.Wait()

		pm.mu.Lock()
		defer pm.mu.Unlock()

		if _, ok := pm.commands[cmd.Process.Pid]; ok {
			delete(pm.commands, pid)
		}

	}(cmd, pid)

	return nil
}
