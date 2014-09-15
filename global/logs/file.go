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

// 已知问题:
// 1.无缓存设计，对系统性能有一定影响，若使用缓存需要解决当系统异常终止时缓冲区的内容无法写入磁盘的问题
// 2.与控制台上的日志不同步(已解决)
// 3.未管理日志文件大小

package logs

import (
    "fmt"
    //"bufio"
    "os"
    "strconv"
    "log"
    "macheart/global/conf"
)

type FileLogger struct {
    logger *log.Logger
    level int
}

// define default value
var (
    defaultFileLoggerLevel = levelInfo
    defaultLogFilePath = "./macheart.log"
)

// create FileLogger returning as LoggerInterface
func NewFileLogger()(LoggerInterface, error){

    filelogger := new(FileLogger)

    // set logger object
    logpath, ok := conf.Get["LogFilePath"]
    if !ok {
        logpath = defaultLogFilePath // set default LogFilePath
    }
    logfile, err:= os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if  err != nil {
        return nil, err
    }
    //logfilebuf := bufio.NewWriter(logfile)
    logfile.WriteString("================macheart log================\n") // write log start flag
    filelogger.logger = log.New(logfile, "", log.LstdFlags)

    // set level
    level, err := strconv.Atoi(conf.Get["FileLoggerLevel"])
    if err != nil || level < levelFatal || level > levelDebug {
        level = defaultFileLoggerLevel // set default level
    }
    filelogger.level = level

    return filelogger, nil
}

// log Fatal level message to file
func (fl *FileLogger)Fatal(format string, v ...interface {}) {

    if levelFatal <= fl.level {
        fl.logger.Printf(format, v...)
    }

}

// log Error level message to file
func (fl *FileLogger)Error(format string, v ...interface {}) {

    if levelError <= fl.level {
        fl.logger.Printf(format, v...)
    }

}

// log Warn level message to file
func (fl *FileLogger)Warn(format string, v ...interface {}) {

    if levelWarn <= fl.level {
        fl.logger.Printf(format, v...)
    }

}

// log Info level message to file
func (fl *FileLogger)Info(format string, v ...interface {}) {

    if levelInfo <= fl.level {
        fl.logger.Printf(format, v...)
    }

}

// log Debug level message to file
func (fl *FileLogger)Debug(format string, v ...interface {}) {

    if levelDebug <= fl.level {
        fl.logger.Printf(format, v...)
    }

}

// register FileLogger in init() function of logs package
func init() {

    isEnable, ok := conf.Get["EnableFileLogger"]
    if isEnable != "0" && ok {
        fileLogger, err:= NewFileLogger()
        if err != nil {
            fmt.Printf("init File Logger fail: %v\n", err)
            return
        }
        register(fileLogger)
    }

}
