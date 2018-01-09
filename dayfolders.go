// dayfolders is a command line tool that creates daily folders to store files in a selectable period of time.
// Copyright (C) 2017-2018  Danilo Cicerone

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

type Locale string

type Format uint8

const (
	short Format = iota
	long
	noname
)

const (
	Sunday time.Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

const (
	January time.Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

var DayNames = map[Locale]map[Format]map[time.Weekday]string{
	"en_US": {
		short: {
			Sunday:    "Sun",
			Monday:    "Mon",
			Tuesday:   "Tue",
			Wednesday: "Wed",
			Thursday:  "Thu",
			Friday:    "Fri",
			Saturday:  "Sat",
		},
		long: {
			Sunday:    "Sunday",
			Monday:    "Monday",
			Tuesday:   "Tuesday",
			Wednesday: "Wednesday",
			Thursday:  "Thursday",
			Friday:    "Friday",
			Saturday:  "Saturday",
		},
	},
	"it_IT": {
		short: {
			Sunday:    "Dom",
			Monday:    "Lun",
			Tuesday:   "Mar",
			Wednesday: "Mer",
			Thursday:  "Gio",
			Friday:    "Ven",
			Saturday:  "Sab",
		},
		long: {
			Sunday:    "Domenica",
			Monday:    "Lunedì",
			Tuesday:   "Martedì",
			Wednesday: "Mercoledì",
			Thursday:  "Giovedì",
			Friday:    "Venerdì",
			Saturday:  "Sabato",
		},
	},
}

var MonthNames = map[Locale]map[Format]map[time.Month]string{
	"en_US": {
		short: {
			January:   "Jan",
			February:  "Feb",
			March:     "Mar",
			April:     "Apr",
			May:       "May",
			June:      "Jun",
			July:      "Jul",
			August:    "Aug",
			September: "Sep",
			October:   "Oct",
			November:  "Nov",
			December:  "Dec",
		},
		long: {
			January:   "January",
			February:  "February",
			March:     "March",
			April:     "April",
			May:       "May",
			June:      "June",
			July:      "July",
			August:    "August",
			September: "September",
			October:   "October",
			November:  "November",
			December:  "December",
		},
	},
	"it_IT": {
		short: {
			January:   "Gen",
			February:  "Feb",
			March:     "Mar",
			April:     "Apr",
			May:       "Mag",
			June:      "Giu",
			July:      "Lug",
			August:    "Aug",
			September: "Set",
			October:   "Ott",
			November:  "Nov",
			December:  "Dic",
		},
		long: {
			January:   "Gennaio",
			February:  "Febbraio",
			March:     "Marzo",
			April:     "Aprile",
			May:       "Maggio",
			June:      "Giugno",
			July:      "Luglio",
			August:    "Agosto",
			September: "Settembre",
			October:   "Ottobre",
			November:  "Novembre",
			December:  "Dicembre",
		},
	},
}

var lc Locale  // Locale Code
var fmn Format // Short or Long month's name
var fdn Format // Short or Long day's name

var yearPtr, fromPtr, toPtr, pathPtr string
var onefPtr, subfPtr, doyPtr, verPtr bool
var daysPtr, dowPtr, domPtr int
var dateStart, dateEnd time.Time
var out io.Writer = ioutil.Discard

const debugModeActivated = false

//------------------------------------------------------------------------------

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage: dayfolders [-year YYYY] [-from YYYY-MM-DD] " +
			"[-to YYYY-MM-DD] [-path /your/target/dir/] [-days 1 to 366] " +
			"[-sub] [-one] [-dow] [-doy] [-ver]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Printf("\nCopyright (C) 2017-2018 Danilo Cicerone.\n" +
			"This is free software; see the source for copying conditions. " +
			"There is NO warranty; not even for MERCHANTABILITY or FITNESS " +
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

	flag.BoolVar(&onefPtr, "one", false, "A long forlder name per "+
		"day(eg: 2017-01-22).")

	flag.BoolVar(&subfPtr, "sub", false, "A forlder for the year, subfolders "+
		"for the months and subfolders for the days(default).")

	flag.IntVar(&daysPtr, "days", 0, "Number of the days(1 to 366) to add or "+
		"subtract(only in case of option '-to').")

	flag.IntVar(&dowPtr, "dow", 0, "0=No day's name added(default), 1=Adds "+
		"short names of the day, 2=Adds names of the day.")

	flag.IntVar(&domPtr, "dom", 0, "0=No month's name added(default), 1=Adds "+
		"short names of the month, 2=Adds names of the month. Doesn't affect"+
		" 'one' option.")

	flag.BoolVar(&doyPtr, "doy", false, "Adds day of the year.")

	flag.BoolVar(&verPtr, "ver", false, "Prints version number.")

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
	fmt.Fprintln(out, "verPtr: ", verPtr)
	fmt.Fprintln(out, "onefPtr: ", onefPtr)
	fmt.Fprintln(out, "subfPtr: ", subfPtr)
	fmt.Fprintln(out, "pathPtr (before Abs): ", pathPtr)
	lc = "it_IT"

	if verPtr {
		fmt.Println("Version 1.1.0")
		os.Exit(0)
	}

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

	if dowPtr < 0 || dowPtr > 2 {
		return errors.New("Ambiguous format for day's name. Goodbye!")
	} else {
		switch dowPtr {
		case 1:
			fdn = short
		case 2:
			fdn = long
		default:
			fdn = noname
		}
	}

	if domPtr < 0 || domPtr > 2 {
		return errors.New("Ambiguous format for month's name. Goodbye!")
	} else {
		switch domPtr {
		case 1:
			fmn = short
		case 2:
			fmn = long
		default:
			fmn = noname
		}
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

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

//------------------------------------------------------------------------------

func createFolders() error {
	var dayOfYear, dayName, monthName string
	dateCursor := dateStart
	yearFormat := "2006"
	monthFormat := "01"
	dayFormat := "02"

	for dateCursor.Before(dateEnd) || dateCursor.Equal(dateEnd) {
		if doyPtr {
			dayOfYear = fmt.Sprintf("%03d", dateCursor.YearDay())
		}

		dayName = fmt.Sprintf("%s", DayNames[lc][fdn][dateCursor.Weekday()])
		monthName = fmt.Sprintf("%s", MonthNames[lc][fmn][dateCursor.Month()])

		if (!onefPtr && !subfPtr) || subfPtr {
			s := []string{dayOfYear, dateCursor.Format(dayFormat), dayName}
			m := []string{dateCursor.Format(monthFormat), monthName}
			err := os.MkdirAll(
				filepath.Join(
					pathPtr,
					dateCursor.Format(yearFormat),
					strings.Join(deleteEmpty(m), "_"),
					strings.Join(deleteEmpty(s), "_"),
				), os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			s := []string{dateCursor.Format("2006-01-02"), dayOfYear, dayName}
			err := os.MkdirAll(filepath.Join(
				pathPtr,
				dateCursor.Format(yearFormat),
				strings.Join(deleteEmpty(s), "_"),
			), os.ModePerm)

			if err != nil {
				return err
			}
		}

		dateCursor = dateCursor.AddDate(0, 0, 1)
	}

	return nil
}
