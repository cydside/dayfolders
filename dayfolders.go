// dayfolders is a command line tool that creates daily folders to store files in a selectable period of time.
// Copyright (C) 2017  Danilo Cicerone

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

//------------------------------------------------------------------------------

package main

//------------------------------------------------------------------------------

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//------------------------------------------------------------------------------

var yearPtr, fromPtr, toPtr, pathPtr string
var onefPtr, subfPtr, dowPtr, doyPtr bool
var daysPtr int
var dateStart, dateEnd time.Time
var out io.Writer = ioutil.Discard

const debugModeActivated = false

//------------------------------------------------------------------------------

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage: dayfolders [-year YYYY] [-from YYYY-MM-DD] " +
			"[-to YYYY-MM-DD] [-path /your/target/dir/] [-days 1 to 366] " +
			"[-sub] [-one] [-dow] [-doy]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Printf("\nCopyright (C) 2017 Danilo Cicerone.\n" +
			"This is free software; see the source for copying conditions. "+
			"There is NO warranty; not even for MERCHANTABILITY or FITNESS "+
			"FOR A PARTICULAR PURPOSE.\n")
	}

	flag.StringVar(&yearPtr, "year", "", "Creates folders for each day of "+
		"the requested year.")

	flag.StringVar(&fromPtr, "from", "", "The starting date, if the final "+
		"date won't be supplied the last day of the year will be the selected.")

	flag.StringVar(&toPtr, "to", "", "The final date, if the starting date "+
		"won't be supplied the first day of the year will be the selected.")

	flag.StringVar(&pathPtr, "path", ".", "Path to dir within the folders "+
		"will be created.")

	flag.BoolVar(&onefPtr, "one", false, "A forlder per day(eg: 2017-01-22).")

	flag.BoolVar(&subfPtr, "sub", false, "A forlder for the year, subfolders "+
		"for the months and subfolders for the days(default).")

	flag.IntVar(&daysPtr, "days", 0, "Number of the days(1 to 366) to add or "+
		"subtract(only in case of option '-to').")

	flag.BoolVar(&dowPtr, "dow", false, "Adds short names of the day.")

	flag.BoolVar(&doyPtr, "doy", false, "Adds day of the year.")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	if debugModeActivated {
		out = os.Stdout
	}

	fmt.Fprintln(out, "yearPtr: ", yearPtr)
	fmt.Fprintln(out, "fromPtr: ", fromPtr)
	fmt.Fprintln(out, "toPtr: ", toPtr)
	fmt.Fprintln(out, "daysPtr: ", daysPtr)
	fmt.Fprintln(out, "dowPtr: ", dowPtr)
	fmt.Fprintln(out, "doyPtr: ", doyPtr)
	fmt.Fprintln(out, "onefPtr: ", onefPtr)
	fmt.Fprintln(out, "subfPtr: ", subfPtr)
	fmt.Fprintln(out, "pathPtr (before Abs): ", pathPtr)

	if err := validateFlags(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintln(out, "pathPtr: ", pathPtr)

	if err := validatePeriod(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintln(out, "dateStart: ", dateStart)
	fmt.Fprintln(out, "dateEnd: ", dateEnd)

	if err := createFolders(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//------------------------------------------------------------------------------

func validateFlags() error {

	var errPeriod error = errors.New("Ambiguous period. Goodbye!")

	if yearPtr == "" && fromPtr == "" && toPtr == "" {
		return errPeriod
	}

	if daysPtr > 0 && fromPtr != "" && toPtr != "" {
		return errPeriod
	}

	if (yearPtr != "" && fromPtr != "") || (yearPtr != "" && toPtr != "") {
		return errPeriod
	}

	if onefPtr && subfPtr {
		return errors.New("Ambiguous style for folders. Goodbye!")
	}

	if pathPtr == "." {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}
		pathPtr = dir
	} else {
		if _, err := os.Stat(pathPtr); os.IsNotExist(err) {
			return err // Dir doesn't exist!
		} else {
			pathPtr, _ = filepath.Abs(filepath.Dir(pathPtr))
		}
	}

	if daysPtr < 0 || daysPtr > 366 {
		return errors.New("Days between 1 and 366. Goodbye!")
	}

	return nil
}

//------------------------------------------------------------------------------

func validatePeriod() error {

	var errDateFormat error = errors.New("Unexpected date format. Goodbye!")

	if yearPtr != "" {
		test, err := time.Parse("2006", yearPtr)
		if err != nil {
			return errDateFormat
		}
		dateStart = test
		if daysPtr > 0 {
			dateEnd = dateStart.AddDate(0, 0, daysPtr-1)
		} else {
			dateEnd = dateStart.AddDate(1, 0, 0).Add(-time.Nanosecond)
		}
	}

	if fromPtr != "" {
		test, err := time.Parse("2006-01-02", fromPtr)
		if err != nil {
			return errDateFormat
		}
		dateStart = test

		if toPtr == "" {
			if daysPtr > 0 {
				dateEnd = dateStart.AddDate(0, 0, daysPtr-1)
			} else {
				dateEnd, _ = time.Parse("2006", strconv.Itoa(dateStart.Year()))
				dateEnd = dateEnd.AddDate(1, 0, 0).Add(-time.Nanosecond)
			}
		}
	}

	if toPtr != "" {
		test, err := time.Parse("2006-01-02", toPtr)
		if err != nil {
			return errDateFormat
		}
		dateEnd = test

		if fromPtr == "" {
			if daysPtr > 0 {
				dateStart = dateEnd.AddDate(0, 0, -daysPtr+1)
			} else {
				dateStart, _ = time.Parse("2006", strconv.Itoa(dateEnd.Year()))
			}
		}
	}

	if dateEnd.Before(dateStart) {
		return errors.New("Starting date is after the final. Goodbye!")
	}

	return nil
}

//------------------------------------------------------------------------------

func createFolders() error {

	dayOfYear := ""
	dateCursor := dateStart
	yearFormat := "2006"
	monthFormat := "01"
	dayFormat := "02"
	shortDayNameFormat := ""

	if dowPtr {
		shortDayNameFormat = " (Mon)"
	}

	for dateCursor.Before(dateEnd) || dateCursor.Equal(dateEnd) {
		if doyPtr {
			dayOfYear = fmt.Sprintf("%03d", dateCursor.YearDay())
		}

		if (!onefPtr && !subfPtr) || subfPtr {
			err := os.MkdirAll(
				filepath.Join(
					pathPtr,
					dateCursor.Format(yearFormat),
					dateCursor.Format(monthFormat),
					strings.TrimSpace(
						dateCursor.Format(dayFormat+
							shortDayNameFormat)+" "+dayOfYear),
				), os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			err := os.Mkdir(filepath.Join(
				pathPtr,
				strings.TrimSpace(
					dateCursor.Format("2006-01-02"+shortDayNameFormat)+
						" "+dayOfYear),
			), os.ModePerm)

			if err != nil {
				return err
			}
		}

		dateCursor = dateCursor.AddDate(0, 0, 1)
	}

	return nil
}
