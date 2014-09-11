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

package main

import (
    "os"
    "net"
    "runtime"
    "strconv"
    "macheart/global"
)

func main() {

    global.Logger.Debug("hehe%s%s" , "hehe", "lijie")
    global.Logger.Debug("hehe")
    global.Logger.Info("niemi")
    global.Logger.Warn("asfdas")

    // set GOMAXPROCS
    cpus, e := strconv.Atoi(global.Configure["CPUs"])
    if e != nil {
        cpus = 1
    }
    runtime.GOMAXPROCS(cpus)

    go global.Logger.Warn("asfdas1")
    go global.Logger.Warn("asfdas2")
    go global.Logger.Warn("asfdas3")
    go global.Logger.Warn("asfdas4")
    go global.Logger.Warn("asfdas5")

    // create Listener
    protocolTypeStr := global.Configure["ProtocolType"]
    listenAddrStr := global.Configure["ListenAddr"]

    listenAddr, e := net.ResolveTCPAddr(protocolTypeStr, listenAddrStr)
    if e != nil {
        global.Logger.Fatal("macheart resolve ListenAddr fail: %s", e.Error())
        os.Exit(1)
        return
    }
    listener, e := net.ListenTCP(protocolTypeStr, listenAddr)
    if e != nil {
        global.Logger.Fatal("macheart create listenner fail: %s", e.Error())
        os.Exit(1)
        return
    }
    defer listener.Close()

    global.Logger.Info("macheart start success!")
    // start listen
    for {
        tcpConn, e := listener.AcceptTCP()
        if e != nil {
            global.Logger.Fatal("macheart accept link of client fail: ", e.Error())
            os.Exit(1)
            return
        }

        // 客户端用户验证
        ok, userID := verifyUser(tcpConn)
        if !ok {
            tcpConn.Close()
            continue
        }


        global.Logger.Info("user %s(%v) is logged\n", userID, tcpConn.RemoteAddr())
        // 处理客户端请求
        go func(conn *net.TCPConn, userID string) {
            //TODO
        }(tcpConn, userID)
    }


}

// 验证用户
func verifyUser(conn *net.TCPConn) (result bool, userID string) {
    return true, "test"//TODO
}

