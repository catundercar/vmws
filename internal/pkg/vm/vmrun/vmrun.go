package vmrun

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type VMRun struct {
}

// List returns all running VMs path.
func (vmrun *VMRun) List(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "vmrun", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %s %w", cmd.String(), output, err)
	}

	var paths []string
	s := bufio.NewScanner(bytes.NewBuffer(output))
	for s.Scan() {
		paths = append(paths, s.Text())
	}
	if len(paths) != 0 {
		paths = paths[1:]
	}
	return paths, nil
}

// Start a vm by path.
func (vmrun *VMRun) Start(ctx context.Context, path string) error {
	sh := fmt.Sprintf("vmrun start %s nogui", path)
	// bug: 命令不会退出。
	cmd := exec.CommandContext(ctx, "sh", "-c", sh)
	cmd.WaitDelay = 10 * time.Second
	b, err := cmd.CombinedOutput()
	if err != nil && !errors.Is(err, exec.ErrWaitDelay) {
		return fmt.Errorf("exec command: %s %w, output: %s", sh, err, string(b))
	}
	log.Println(string(b))
	return nil
}

func (vmrun *VMRun) ctx(ctx context.Context) (context.Context, context.CancelFunc) {
	_, ok := ctx.Deadline()
	if !ok {
		ctxT, cancel := context.WithTimeout(ctx, 30*time.Second)
		return ctxT, cancel
	}
	return ctx, func() {}
}
