package docker

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/daniel1302/vega-asistant/types"
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

func TemplatePostgresqlDockerCompose(
	credentials types.SQLCredentials,
	homePath string,
) (string, error) {
	tmpl := template.Must(template.New("docker-compose.yaml").Parse(postgresqlTemplate))

	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, struct {
		Username string
		DbName   string
		Password string
		Port     int
	}{
		Username: credentials.User,
		DbName:   credentials.DatabaseName,
		Password: credentials.Pass,
		Port:     credentials.Port,
	}); err != nil {
		return "", fmt.Errorf("failed to template run-config.toml: %w", err)
	}

	return buff.String(), nil

	return "", nil
}
