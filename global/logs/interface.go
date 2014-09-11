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

package logs

// define Log Levels
const  (
    levelFatal = iota
    levelError
    levelWarn
    levelInfo
    levelDebug
)

// define Logger Interfaces
type LoggerInterface interface {
    Fatal(format string, v ...interface {})
    Error(format string, v ...interface {})
    Warn(format string, v ...interface {})
    Info(format string, v ...interface {})
    Debug(format string, v ...interface {})
}
