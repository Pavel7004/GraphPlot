/*
Copyright © 2022 Kovalev Pavel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/package main

import (
	"os"

	"github.com/Pavel7004/Common/tracing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Pavel7004/GraphPlot/cmd"
)

func main() {
	closer := tracing.InitDefaultJaeger("GraphPlot")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cmd.Execute()

	defer closer.Close()
}
