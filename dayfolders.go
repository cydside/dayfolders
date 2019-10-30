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

//______________________________________________________________________________

package main

//______________________________________________________________________________

import (
	"fmt"

	"github.com/cydside/dayfolders/cmd"
)

//______________________________________________________________________________

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
