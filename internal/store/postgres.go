package store

import (
	"context"
	"os/exec"
)

func RestoreToDocker(ctx context.Context, dumpFile string) error {
	cmd := exec.CommandContext(ctx,
		"docker", "run", "--rm",
		"-e", "POSTGRES_PASSWORD=test",
		"postgres:15",
	)
	if err := cmd.Start(); err != nil {
		return err
	}

	restore := exec.CommandContext(ctx,
		"docker", "exec", "-i",
		"psql", "-U", "postgres", "-f", dumpFile,
	)

	return restore.Run()
}
