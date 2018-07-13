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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//------------------------------------------------------------------------------

type locale string

type nameAbbr uint8

const (
	short nameAbbr = iota
	long
	noname
)

const (
	sunday time.Weekday = iota
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
)

const (
	january time.Month = 1 + iota
	february
	march
	april
	may
	june
	july
	august
	september
	october
	november
	december
)

var daynames = map[locale]map[nameAbbr]map[time.Weekday]string{
	"en_US": {
		short: {
			sunday:    "Sun",
			monday:    "Mon",
			tuesday:   "Tue",
			wednesday: "Wed",
			thursday:  "Thu",
			friday:    "Fri",
			saturday:  "Sat",
		},
		long: {
			sunday:    "Sunday",
			monday:    "Monday",
			tuesday:   "Tuesday",
			wednesday: "Wednesday",
			thursday:  "Thursday",
			friday:    "Friday",
			saturday:  "Saturday",
		},
	},
	"it_IT": {
		short: {
			sunday:    "Dom",
			monday:    "Lun",
			tuesday:   "Mar",
			wednesday: "Mer",
			thursday:  "Gio",
			friday:    "Ven",
			saturday:  "Sab",
		},
		long: {
			sunday:    "Domenica",
			monday:    "Lunedì",
			tuesday:   "Martedì",
			wednesday: "Mercoledì",
			thursday:  "Giovedì",
			friday:    "Venerdì",
			saturday:  "Sabato",
		},
	},
}

var monthNames = map[locale]map[nameAbbr]map[time.Month]string{
	"en_US": {
		short: {
			january:   "Jan",
			february:  "Feb",
			march:     "Mar",
			april:     "Apr",
			may:       "May",
			june:      "Jun",
			july:      "Jul",
			august:    "Aug",
			september: "Sep",
			october:   "Oct",
			november:  "Nov",
			december:  "Dec",
		},
		long: {
			january:   "January",
			february:  "February",
			march:     "March",
			april:     "April",
			may:       "May",
			june:      "June",
			july:      "July",
			august:    "August",
			september: "September",
			october:   "October",
			november:  "November",
			december:  "December",
		},
	},
	"it_IT": {
		short: {
			january:   "Gen",
			february:  "Feb",
			march:     "Mar",
			april:     "Apr",
			may:       "Mag",
			june:      "Giu",
			july:      "Lug",
			august:    "Aug",
			september: "Set",
			october:   "Ott",
			november:  "Nov",
			december:  "Dic",
		},
		long: {
			january:   "Gennaio",
			february:  "Febbraio",
			march:     "Marzo",
			april:     "Aprile",
			may:       "Maggio",
			june:      "Giugno",
			july:      "Luglio",
			august:    "Agosto",
			september: "Settembre",
			october:   "Ottobre",
			november:  "Novembre",
			december:  "Dicembre",
		},
	},
}

var lc locale    // locale Code
var fmn nameAbbr // Short or Long month's name
var fdn nameAbbr // Short or Long day's name

var yearPtr, fromPtr, toPtr, pathPtr, prefPtr, suffPtr string
var onefPtr, subfPtr, doyPtr, verPtr bool
var daysPtr, dowPtr, domPtr int
var dateStart, dateEnd time.Time
var out = ioutil.Discard

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

	flag.StringVar(&suffPtr, "suffix", "", "Add suffix to day folder's name.")

	flag.StringVar(&prefPtr, "prefix", "", "Add prefix to day folder's name.")

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
	fmt.Fprintln(out, "prefPtr: ", prefPtr)
	fmt.Fprintln(out, "suffPtr: ", suffPtr)
	fmt.Fprintln(out, "daysPtr: ", daysPtr)
	fmt.Fprintln(out, "dowPtr: ", dowPtr)
	fmt.Fprintln(out, "doyPtr: ", doyPtr)
	fmt.Fprintln(out, "verPtr: ", verPtr)
	fmt.Fprintln(out, "onefPtr: ", onefPtr)
	fmt.Fprintln(out, "subfPtr: ", subfPtr)
	fmt.Fprintln(out, "pathPtr (before Abs): ", pathPtr)
	lc = "it_IT"

	if verPtr {
		fmt.Println("Version 1.2.0")
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
	var errPeriod = errors.New("ambiguous period, goodbye")
	var PathSeparator = fmt.Sprintf("%c", os.PathSeparator)

	if yearPtr == "" && fromPtr == "" && toPtr == "" {
		return errPeriod
	}

	if daysPtr > 0 && fromPtr != "" && toPtr != "" {
		return errPeriod
	}

	if strings.Index(suffPtr, PathSeparator) > 0 {
		return errors.New("suffix text contains os path separator, goodbye")
	}

	if strings.Index(prefPtr, PathSeparator) > 0 {
		return errors.New("prefix text contains os path separator, goodbye")
	}

	if (yearPtr != "" && fromPtr != "") || (yearPtr != "" && toPtr != "") {
		return errPeriod
	}

	if onefPtr && subfPtr {
		return errors.New("ambiguous style for folders, goodbye")
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
		}

		pathPtr, _ = filepath.Abs(filepath.Dir(pathPtr))
	}

	if daysPtr < 0 || daysPtr > 366 {
		return errors.New("days between 1 and 366, goodbye")
	}

	if dowPtr < 0 || dowPtr > 2 {
		return errors.New("ambiguous format for day's name, goodbye")
	}

	switch dowPtr {
	case 1:
		fdn = short
	case 2:
		fdn = long
	default:
		fdn = noname
	}

	if domPtr < 0 || domPtr > 2 {
		return errors.New("ambiguous format for month's name, goodbye")
	}

	switch domPtr {
	case 1:
		fmn = short
	case 2:
		fmn = long
	default:
		fmn = noname
	}

	return nil
}

//------------------------------------------------------------------------------

func validatePeriod() error {
	var errDF = errors.New("unexpected date format, goodbye")

	if yearPtr != "" {
		test, err := time.Parse("2006", yearPtr)
		if err != nil {
			return errDF
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
			return errDF
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
			return errDF
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
		return errors.New("starting date is after the final, goodbye")
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
	var y, d, m string
	p := dateStart
	yf := "2006"
	mf := "01"
	df := "02"

	for p.Before(dateEnd) || p.Equal(dateEnd) {
		if doyPtr {
			y = fmt.Sprintf("%03d", p.YearDay())
		}

		d = fmt.Sprintf("%s", daynames[lc][fdn][p.Weekday()])
		m = fmt.Sprintf("%s", monthNames[lc][fmn][p.Month()])

		if (!onefPtr && !subfPtr) || subfPtr {
			s := []string{prefPtr, y, p.Format(df), d, suffPtr}
			n := []string{p.Format(mf), m}
			err := os.MkdirAll(
				filepath.Join(
					pathPtr,
					p.Format(yf),
					strings.Join(deleteEmpty(n), "_"),
					strings.Join(deleteEmpty(s), "_"),
				), os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			s := []string{prefPtr, p.Format("2006-01-02"), y, d, suffPtr}
			err := os.MkdirAll(filepath.Join(
				pathPtr,
				p.Format(yf),
				strings.Join(deleteEmpty(s), "_"),
			), os.ModePerm)

			if err != nil {
				return err
			}
		}

		p = p.AddDate(0, 0, 1)
	}

	return nil
}
