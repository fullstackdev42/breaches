package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type Person struct {
	ID1, ID2, FirstName, LastName, Gender, BirthPlace, CurrentPlace, Job, Date string
}

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the data in a sortable table",
	Long: `This command reads the data from the specified file and displays it in a sortable table.
You can navigate through the table using the next and back buttons.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open("data/Canada.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		var people []Person
		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) >= 10 {
				people = append(people, Person{
					ID1:          parts[0],
					ID2:          parts[1],
					FirstName:    parts[2],
					LastName:     parts[3],
					Gender:       parts[4],
					BirthPlace:   parts[5],
					CurrentPlace: parts[6],
					Job:          parts[8],
					Date:         parts[9],
				})
			} else {
				fmt.Printf("Line %d does not have at least 10 parts: %s\n", lineNumber, line)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID1", "ID2", "First Name", "Last Name", "Gender", "Birth Place", "Current Place", "Job", "Date"})

		for _, person := range people {
			table.Append([]string{person.ID1, person.ID2, person.FirstName, person.LastName, person.Gender, person.BirthPlace, person.CurrentPlace, person.Job, person.Date})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
