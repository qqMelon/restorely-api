package restore

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func TestPostgresRestoreDocker(ctx context.Context, dumpPath string) error {
	container := fmt.Sprintf("restorely-test-%d", time.Now().Unix())

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	run := exec.CommandContext(ctx,
		"docker", "run", "-d",
		"--name", container,
		"-e", "POSTGRES_PASSWORD=restorely",
		"-e", "POSTGRES_DB=restoredb",
		"postgres:15",
	)

	if err := run.Run(); err != nil {
		return fmt.Errorf("docker run failed: %w", err)
	}

	defer exec.Command("docker", "rm", "-f", container).Run()

	// Wait for PostgreSQL to be ready
	if err := waitForPostgres(ctx, container); err != nil {
		return err
	}

	// Copy dump
	copy := exec.CommandContext(ctx,
		"docker", "cp",
		dumpPath,
		container+":/backup.dump",
	)

	if err := copy.Run(); err != nil {
		return fmt.Errorf("docker cp failed: %w", err)
	}

	// Restore
	restore := exec.CommandContext(ctx,
		"docker", "exec", container,
		"pg_restore",
		"-U", "postgres",
		"--no-owner",
		"--no-acl",
		"-d", "restoredb",
		"/backup.dump",
	)

	if err := restore.Run(); err != nil {
		return fmt.Errorf("pg_restore failed: %w", err)
	}

	// Basic validation
	if err := validateRestore(ctx, container); err != nil {
		return err
	}

	return nil
}
func waitForPostgres(ctx context.Context, container string) error {
	for i := 0; i < 10; i++ {
		cmd := exec.CommandContext(ctx, "docker", "exec", container, "pg_isready", "-U", "postgres")
		if err := cmd.Run(); err == nil {
			return nil
		}

		time.Sleep(3 * time.Second)
	}

	return fmt.Errorf("postgres not ready")
}

func validateRestore (ctx context.Context, container string) error {
  cmd := exec.CommandContext(ctx,
		"docker", "exec", container,
		"psql",
		"-U", "postgres",
		"-d", "restoredb",
		"-c", "SELECT count(*) FROM information_schema.tables",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if string(out) == "" {
		return fmt.Errorf("empty validation output")
	}

	return nil
}
