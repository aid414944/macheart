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

import (
     "sync"
)

type Logger struct {
    sync.Mutex
    subloggers []LoggerInterface
}

var subloggers = make([]LoggerInterface, 0, 10)

// create Logger
func NewLogger() *Logger {
    lg := new(Logger)
    lg.subloggers = subloggers
    return lg
}

// log Fatal level message
func (lgr *Logger)Fatal(format string, v ...interface {}) {
    lgr.Lock()
    defer lgr.Unlock()
    for _, slg := range lgr.subloggers {
        slg.Fatal("[Fatal]" + format, v...)
    }
}

// log Error level message
func (lgr *Logger)Error(format string, v ...interface {}) {
    lgr.Lock()
    defer lgr.Unlock()
    for _, slg := range lgr.subloggers {
        slg.Error("[Error]" + format, v...)
    }
}

// log Warn level message
func (lgr *Logger)Warn(format string, v ...interface {}) {
    lgr.Lock()
    defer lgr.Unlock()
    for _, slg := range lgr.subloggers {
        slg.Warn("[Warn]" + format, v...)
    }
}

// log Info level message
func (lgr *Logger)Info(format string, v ...interface {}) {
    lgr.Lock()
    defer lgr.Unlock()
    for _, slg := range lgr.subloggers {
        slg.Info("[Info]" + format, v...)
    }
}

// log Debug level message
func (lgr *Logger)Debug(format string, v ...interface {}) {
    lgr.Lock()
    defer lgr.Unlock()
    for _, slg := range lgr.subloggers {
        slg.Debug("[Debug]" + format, v...)
    }
}

// register sublogger
func register(logger LoggerInterface) {
    subloggers = append(subloggers, logger)
}
