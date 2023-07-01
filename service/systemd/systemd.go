package systemd

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"runtime"

	"go.uber.org/zap"

	"github.com/daniel1302/vega-assistant/utils"
)

const systemdTemplate = `[Unit]
Description=vegavisor
Documentation=https://github.com/vegaprotocol/vega
After=network.target network-online.target
Requires=network-online.target

[Service]
User={{.User}}
Group={{.Group}}
ExecStart="{{.VisorHome}}/visor" run --home "{{.VisorHome}}"
TimeoutStopSec=10s
LimitNOFILE=1048576
LimitNPROC=512
PrivateTmp=false
ProtectSystem=full
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target`

const serviceFilePath = "/lib/systemd/system/vegavisor.service"

func PrepareSystemd(logger *zap.SugaredLogger, visorHome string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("systemd supported only on Linux")
	}

	currentUser, err := utils.Whoami()
	if err != nil {
		return fmt.Errorf("failed to get current user name: %w", err)
	}

	ownerUser, ownerGroup, err := utils.GetFileOwner(visorHome)
	if err != nil {
		return fmt.Errorf("failed to describe owner for %s: %w", visorHome, err)
	}

	systemdServiceContent, err := templateSystemdService(visorHome, ownerUser, ownerGroup)
	if err != nil {
		return fmt.Errorf("failed to template systemd service: %w", err)
	}

	if currentUser != "root" || utils.IsWSL() {
		fmt.Println(systemdServiceContent)
		return nil
	}

	logger.Infof("Updating content of the service file in %s", serviceFilePath)
	if err := os.WriteFile(serviceFilePath, []byte(systemdServiceContent), os.ModePerm); err != nil {
		return fmt.Errorf("failed to update %s file: %w", serviceFilePath, err)
	}

	logger.Info("Calling systemctl daemon-reload")
	if _, err := utils.ExecuteBinary("systemctl", []string{"daemon-reload"}, nil); err != nil {
		return fmt.Errorf("failed to call systemctl daemon-reload: %w", err)
	}
	logger.Info("Daemons reloaded")
	return nil
}

func templateSystemdService(visorHome, username, groupname string) (string, error) {
	tmpl := template.Must(template.New("vegavisor.service").Parse(systemdTemplate))

	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, struct {
		VisorHome string
		User      string
		Group     string
	}{
		VisorHome: visorHome,
		User:      username,
		Group:     groupname,
	}); err != nil {
		return "", fmt.Errorf("failed to template vegavisor.service: %w", err)
	}

	return buff.String(), nil
}
