// dayfolders is a command line tool that creates daily folders to store files for a selectable period of time.
// Copyright (C) 2017-2019  Danilo Cicerone

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

//______________________________________________________________________________

import (
	"fmt"
	"os"

	"github.com/cydside/dailyname"
	"github.com/spf13/cobra"
)

//______________________________________________________________________________

var (
	rootCmd = &cobra.Command{
		Use:   "dayfolders",
		Short: "Creates daily folders for a selectable period of time",
		Long:  `A tool that creates daily folders for a selectable period of time`,
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Printf("dayfolders: %v\n", usrReq)
			arrstr, err := dailyname.GetDailyNames(&usrReq)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			for _, el := range arrstr {
				err := os.MkdirAll(el, os.ModePerm)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					os.Exit(2)
				}
			}
		},
	}
)

//______________________________________________________________________________

var usrReq dailyname.UserReq

//______________________________________________________________________________

func init() {
	rootCmd.Flags().StringVarP(&usrReq.Lang, "lang", "l", "en_US", "Language selection: 'it_IT' for italian and for 'en_US' english")
	rootCmd.Flags().StringVarP(&usrReq.DateFrom, "from", "f", "", "The starting date in format YYYY-MM-DD(eg.: 2019-11-08 or 2019-11 or 2019) on missing day or month will be selected the first avaiable. If the final date won't be supplied the last day of the year or month will be the selected (required)")
	rootCmd.MarkFlagRequired("from")
	rootCmd.Flags().StringVarP(&usrReq.DateTo, "to", "t", "", "The final date in format YYYY-MM-DD(eg.: 2020-01-16 or 2020-01) on missing day will be selected the last avaiable")
	rootCmd.Flags().StringVarP(&usrReq.Suffix, "suffix", "u", "", "Add suffix to day folder's name")
	rootCmd.Flags().StringVarP(&usrReq.Prefix, "prefix", "r", "", "Add suffix to day folder's name")
	rootCmd.Flags().BoolVarP(&usrReq.LoneOrSub, "subf", "s", false, "Creates a forlder for the year, subfolders for the months and subfolders for the days otherwise a long forlder name per day, eg: 2017-01-22 (default)")
	rootCmd.Flags().IntVarP(&usrReq.Duration, "duration", "d", 0, "Duration in number of the days(1 to 366) if no end date supplied")
	rootCmd.Flags().StringVarP(&usrReq.Content, "content", "c", "", "Last names to add to daily name separated by comma, eg.: 'DONE,TODO' will do 2019-11-08\\DONE and 2019-11-08\\TODO")
	rootCmd.Flags().IntVarP(&usrReq.DoW, "dow", "w", -1, "Day's name added: Monday(1=long format), Mon(0=short format)")
	rootCmd.Flags().IntVarP(&usrReq.DoM, "dom", "m", -1, "Month's name added: Jenuary(1=long format), Jen(0=short format)")
	rootCmd.Flags().BoolVarP(&usrReq.DoY, "doy", "j", false, "Julian date added: 001 to 365/366")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(languageCmd)
}

//______________________________________________________________________________

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
