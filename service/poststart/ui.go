package poststart

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func printSummary(settings ServiceSettings) {
	fmt.Println("\n Summary:\n")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Parameter", "Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow("Vega Home", settings.VegaHome)
	tbl.AddRow("Tendermint Home", settings.TendermintHome)
	tbl.Print()
	fmt.Println("")
}
