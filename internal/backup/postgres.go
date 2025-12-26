package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type PostgresConfig struct {
	DSN string
}

func BackupPostgres(ctx context.Context, cfg PostgresConfig, outputPath string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"pg_dump",
		"--format=custom",
		"--no-owner",
		"--no-aci",
		cfg.DSN,
	)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	return nil
}
