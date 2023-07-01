package postgresql

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func PrintInstructions(homePath string) {
	fmt.Printf(`

    Your setup is ready. Now you have to start postgreSQL with the following commands:

    cd %s;
    docker-compose up -d;`, homePath)
	fmt.Println("")
}

func printSummary(settings GeneratorSettings) {
	fmt.Println("\n Summary:\n")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Parameter", "Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow("Home", settings.Home)
	tbl.AddRow("SQL Port", settings.PostgresqlPort)
	tbl.AddRow("SQL User", settings.PostgresqlUsername)
	tbl.AddRow(
		"SQL Password",
		fmt.Sprintf(
			"%c***%c",
			settings.PostgresqlPassword[0],
			settings.PostgresqlPassword[len(settings.PostgresqlPassword)-1],
		),
	)
	tbl.AddRow("SQL Database Name", settings.PostgresqlDatabase)

	tbl.Print()
	fmt.Println("")
}
