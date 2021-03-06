/*
macheart Copyright (C) 2014  aid414944

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package global

import (
    "macheart/conf"
    "macheart/logs"
    "path/filepath"
    "os/exec"
)


var (
    // init Configure
    Configure = conf.Get
    // init Logger
    Logger = logs.NewLogger()
)

// the function is used to fork macheart oneself
func ForkOneself(args []string) error {
    filePath, _ := filepath.Abs(args[0])
    cmd := exec.Command(filePath, args[1:]...)
    return cmd.Start()
}
