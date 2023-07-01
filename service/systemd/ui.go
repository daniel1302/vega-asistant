package systemd

import (
	"fmt"

	"github.com/daniel1302/vega-assistant/utils"
)

func PrintInstructions() {
	currentUser, _ := utils.Whoami()
	if currentUser == "root" && !utils.IsWSL() {
		fmt.Println(`
      Systemd service installed. You can use following command to start your node:
      
        sudo systemctl start vegavisor

      You can see the node logs with the following command:

        sudo journalctl -u vegavisor -n 1000 -f`)

		return
	}

	fmt.Println(`
    You MUST manually install the systemd service. To do it:
      1. Create the '/lib/systemd/system/vegavisor.service' file
      2. Put the above content in the created file
      3. Call the 'sudo systemctl daemon-reload' command
    
    You can use the following command to start the node:

     sudo systemctl start vegavisor

    You can see the node logs with the following command:

      sudo journalctl -u vegavisor -n 1000 -f`)
}
