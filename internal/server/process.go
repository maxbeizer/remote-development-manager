package server

import (
	"bytes"
	"context"
	"os/exec"
	"sync"
)

type processManager struct {
	commands map[int]*exec.Cmd
	mu       sync.Mutex
}

func (pm *processManager) RunningProcesses() []*exec.Cmd {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	commands := make([]*exec.Cmd, 0)

	for _, command := range pm.commands {
		commands = append(commands, command)
	}

	return commands
}

func (pm *processManager) AddPid(cmd *exec.Cmd, pid int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.commands[pid] = cmd
}

func (pm *processManager) RemovePid(pid int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.commands, pid)
}

func (pm *processManager) RunNow(ctx context.Context, name string, path string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, path, args...)

	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Start()

	if err != nil {
		return nil, err
	}

	pid := cmd.Process.Pid

	pm.AddPid(cmd, pid)
	defer pm.RemovePid(pid)

	err = cmd.Wait()

	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func (pm *processManager) RunInBackground(ctx context.Context, name string, path string, args ...string) error {
	cmd := exec.CommandContext(ctx, path, args...)

	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Start()

	if err != nil {
		return err
	}

	pm.AddPid(cmd, cmd.Process.Pid)

	go func() {
		defer pm.RemovePid(cmd.Process.Pid)
		cmd.Wait()
	}()

	return nil
}
