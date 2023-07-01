package vegacmd

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/daniel1302/vega-asistant/utils"
)

const VisorRunConfigTemplate = `name = "{{.Version}}"

[vega]
  [vega.binary]
    path = "vega"
    args = ["start", "--home", "{{.VegaHome}}", "--tendermint-home", "{{.TendermintHome}}"]
  [vega.rpc]
    socketPath = "/tmp/vega.sock"
    httpPath = "/rpc"

[data_node]
  [data_node.binary]
    path = "vega"
    args = ["datanode", "start", "--home", "{{.VegaHome}}"]`

func InitVisor(binaryPath, visorHome string) error {
	_, err := utils.ExecuteBinary(binaryPath, []string{"init", "--home", visorHome}, nil)
	if err != nil {
		return fmt.Errorf("failed to init vegavisor: %w", err)
	}

	return nil
}

func TemplateVisorRunConfig(version, vegaHome, tendermintHome string) (string, error) {
	tmpl := template.Must(template.New("run-config.toml").Parse(VisorRunConfigTemplate))
	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, struct {
		Version        string
		VegaHome       string
		TendermintHome string
	}{
		Version:        version,
		VegaHome:       vegaHome,
		TendermintHome: tendermintHome,
	}); err != nil {
		return "", fmt.Errorf("failed to template run-config.toml: %w", err)
	}

	return buff.String(), nil
}
