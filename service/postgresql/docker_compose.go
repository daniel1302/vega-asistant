package postgresql

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const postgresqlTemplate = `version: '3.1'

services:
  db:
    image: timescale/timescaledb:2.8.0-pg14
    restart: always
    environment:
      POSTGRES_USER: {{.Username}}
      POSTGRES_DB: {{.DbName}}
      POSTGRES_PASSWORD: {{.Password}}
    command: [
      "postgres",
      "-c", "max_connections=50",
      "-c", "log_destination=stderr",
      "-c", "work_mem=5MB",
      "-c", "huge_pages=off",
      "-c", "shared_memory_type=sysv",
      "-c", "dynamic_shared_memory_type=sysv",
      "-c", "shared_buffers=2GB",
      "-c", "temp_buffers=5MB",
    ]
    ports:
      - {{.Port}}:5432
    volumes: 
      - pgdata:/var/lib/postgresql/data

  # Adminer can be used for debugging. We recommend dissabling it for production
  adminer:
    image: adminer
    restart: always
    ports:
      - 8082:8080


volumes:
  pgdata:
    driver: local`

func PrepareDockerComposeFile(logger *zap.SugaredLogger, settings GeneratorSettings) error {
	logger.Info("Templating docker-compose.yaml file")
	composerContent, err := templatePostgresqlDockerCompose(
		settings.PostgresqlUsername,
		settings.PostgresqlDatabase,
		settings.PostgresqlPassword,
		settings.PostgresqlPort,
	)
	if err != nil {
		return fmt.Errorf("failed to template docker-compose.yaml file: %w", err)
	}
	logger.Info("Content for docker-compose.yaml generated")

	logger.Infof("Creating home for docker-compose.yaml(%s)", settings.Home)
	if err := os.MkdirAll(settings.Home, os.ModePerm); err != nil {
		return fmt.Errorf(
			"failed to create home dir for postgresql docker-compose(%s): %w",
			settings.Home,
			err,
		)
	}
	logger.Info("Home directory created")

	dockerComposeFilePath := filepath.Join(settings.Home, "docker-compose.yaml")
	logger.Infof("Writing docker-compose file to %s", dockerComposeFilePath)

	if err := os.WriteFile(dockerComposeFilePath, []byte(composerContent), os.ModePerm); err != nil {
		return fmt.Errorf(
			"failed to write docker compose file to %s: %w",
			dockerComposeFilePath,
			err,
		)
	}
	logger.Info("docker-compose.yaml created")
	return nil
}

func templatePostgresqlDockerCompose(
	username, dbName, pass string,
	port int,
) (string, error) {
	tmpl := template.Must(template.New("docker-compose.yaml").Parse(postgresqlTemplate))

	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, struct {
		Username string
		DbName   string
		Password string
		Port     int
	}{
		Username: username,
		DbName:   dbName,
		Password: pass,
		Port:     port,
	}); err != nil {
		return "", fmt.Errorf("failed to template run-config.toml: %w", err)
	}

	return buff.String(), nil

	return "", nil
}
