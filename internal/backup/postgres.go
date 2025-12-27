package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type DatabaseConfig struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName string `json:"db_name"`
}

func RunPostgresBackup(ctx context.Context, d DatabaseConfig, output string) error {
	cmd := exec.CommandContext(ctx,
		"pg_dump",
		"-h", d.Host,
		"-p", fmt.Sprintf("%d", d.Port),
		"-U", d.Username,
		"-d", d.DBName,
		"-f", output,
	)

	cmd.Env = append(os.Environ(),
		"PGPASSWORD="+d.Password,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
