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
    "log"
    "os"
    "strconv"
    "macheart/conf"
)

// ConsoleLogger implements LoggerInterface and writes messages to console
type ConsoleLogger struct {
    logger *log.Logger
    level int
}

// define default value
var (
    defaultConsoleLoggerLevel = levelInfo
)

// create ConsoleLogger returning as LoggerInterface
func NewConsoleLogger() LoggerInterface {

    cl := new(ConsoleLogger)

    // set logger object
    cl.logger = log.New(os.Stdout, "", log.LstdFlags)


    // set level
    level, err := strconv.Atoi(conf.Get["ConsoleLoggerLevel"])
    if err != nil || level < levelFatal || level > levelDebug {
        level = defaultConsoleLoggerLevel // set default level
    }
    cl.level = level

    return cl
}

// log Fatal level message to console
func (cl *ConsoleLogger)Fatal(format string, v ...interface {}) {

    if levelFatal <= cl.level {
        cl.logger.Printf(format, v...)
    }

}

// log Error level message to console
func (cl *ConsoleLogger)Error(format string, v ...interface {}) {

    if levelError <= cl.level {
        cl.logger.Printf(format, v...)
    }

}

// log Warn level message to console
func (cl *ConsoleLogger)Warn(format string, v ...interface {}) {

    if levelWarn <= cl.level {
        cl.logger.Printf(format, v...)
    }

}

// log Info level message to console
func (cl *ConsoleLogger)Info(format string, v ...interface {}) {

    if levelInfo <= cl.level {
        cl.logger.Printf(format, v...)
    }

}

// log Debug level message to console
func (cl *ConsoleLogger)Debug(format string, v ...interface {}) {

    if levelDebug <= cl.level {
        cl.logger.Printf(format, v...)
    }

}

// register ConsoleLogger in init() function of logs package
func init() {

    isEnable, ok := conf.Get["EnableConsoleLogger"]
    if isEnable != "0" && ok {
        register(NewConsoleLogger())
    }

}
