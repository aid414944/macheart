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

package conf

import (
    "fmt"
    "os"
    "strings"
    "io"
    "bufio"
)

// init Configure Data, the Get variable is used to other subsystem of global package
var Get = readConfigFile("./macheart.conf")

// read configure file to create Configure Data
func readConfigFile(filePath string) (confMap map[string]string) {

    f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
    if err != nil {
        fmt.Printf("%v\n", err)
        return nil
    }
    defer f.Close()

    m := make(map[string]string)
    fileReader := bufio.NewReader(f)
    for {
        lineByte, _, e := fileReader.ReadLine()

        if e == io.EOF {
            break
        }

        lineStr := strings.TrimSpace(string(lineByte))
        if lineStr == "" {
            continue
        }
        if lineStr[0] == '#' {
            continue
        }

        keyvalueStr := strings.Split(lineStr, "=")
        keyStr := strings.TrimSpace(keyvalueStr[0])
        valueStr := strings.TrimSpace(keyvalueStr[1])

        m[keyStr] = valueStr
    }

    return m
}
