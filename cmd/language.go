package cmd

//______________________________________________________________________________

import (
	"github.com/spf13/cobra"
)

//______________________________________________________________________________

var languageCmd = &cobra.Command{
	Use:   "language",
	Short: "Used to add a new language set for days' name and months' name",
	Long: `Used to add a new language set for days' name and months' name
		reading from a json file as follow:
		{
			DaysName: {
				"en_US": {
					Short: {
						0: "Sun",
						1: "Mon",
						2: "Tue",
						3: "Wed",
						4: "Thu",
						5: "Fri",
						6: "Sat",
					},
					Long: {
						0: "Sunday",
						1: "Monday",
						2: "Tuesday",
						3: "Wednesday",
						4: "Thursday",
						5: "Friday",
						6: "Saturday",
					},
				},
			},
			MonthsName: {
				"en_US": {
					Short: {
						1:  "Jan",
						2:  "Feb",
						3:  "Mar",
						4:  "Apr",
						5:  "May",
						6:  "Jun",
						7:  "Jul",
						8:  "Aug",
						9:  "Sep",
						10: "Oct",
						11: "Nov",
						12: "Dec",
					},
					Long: {
						1:  "January",
						2:  "February",
						3:  "March",
						4:  "April",
						5:  "May",
						6:  "June",
						7:  "July",
						8:  "August",
						9:  "September",
						10: "October",
						11: "November",
						12: "December",
					},
				}
			}
		}`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("dayfolders args: %v\n", args)
		// fmt.Println("dayfolders a CLI that creates daily folders v1.0")
	},
}
