package process

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"sync"
)

var ErrNoProcess = errors.New("No process with that PID found")

type Manager struct {
	commands map[int]*exec.Cmd
	mu       sync.Mutex
}

func NewManager() *Manager {
	return &Manager{commands: map[int]*exec.Cmd{}}
}

func (m *Manager) RunningProcesses() []*exec.Cmd {
	m.mu.Lock()
	defer m.mu.Unlock()

	commands := make([]*exec.Cmd, 0)

	for _, command := range m.commands {
		commands = append(commands, command)
	}

	return commands
}

func (m *Manager) AddPid(cmd *exec.Cmd, pid int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.commands[pid] = cmd
}

func (m *Manager) RemovePid(pid int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.commands, pid)
}

func (m *Manager) Kill(pid int) error {
	commands := m.RunningProcesses()

	// this is really ineffecient, but probably not a big deal since it's
	// unlikely that RDM will manage a significant number of processes.
	for _, command := range commands {
		if command.Process.Pid == pid {
			return command.Process.Kill()
		}
	}

	return ErrNoProcess
}

func (m *Manager) RunNow(ctx context.Context, name string, path string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, path, args...)

	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Start()

	if err != nil {
		return nil, err
	}

	pid := cmd.Process.Pid

	m.AddPid(cmd, pid)
	defer m.RemovePid(pid)

	err = cmd.Wait()

	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func (m *Manager) RunInBackground(ctx context.Context, name string, path string, args ...string) error {
	cmd := exec.CommandContext(ctx, path, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdin.Close()

	err = cmd.Start()

	if err != nil {
		return err
	}

	m.AddPid(cmd, cmd.Process.Pid)

	go func() {
		defer m.RemovePid(cmd.Process.Pid)
		cmd.Wait()
	}()

	return nil
}
