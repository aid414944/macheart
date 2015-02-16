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
	"fmt"
	"os"
	//"net/rpc"
	"net"
	"runtime"
	"strconv"
	//"path/filepath"
	"macheart/global"
	//"os/exec"
	rpclient "macheart/rpc/client"
	rpcserver "macheart/rpc/server"
)

var helpStr = `Usage: macheart {start|restart|top|--help|--version}`
var versionStr = `Macheart 0.01 by aid414944
Email: aid414944@gmail.com`

func main() {

	if len(os.Args) > 1 {

		// 处理普通参数
		switch os.Args[1] {
		case "--help":
			fmt.Println(helpStr)
			return
		case "--version":
			fmt.Println(versionStr)
			return
		}

		// 处理特殊参数
		cmdHandler := rpclient.NewCmdHandler()
		err := cmdHandler.LinkServer(global.Configure["RpcNetwork"], ":"+global.Configure["RpcPort"])
		if err == nil {
			cmdHandler.Exec(os.Args)
			cmdHandler.Close()
			return
		}

		// 过滤此时无效的命令
		if os.Args[1] != "start" {
			fmt.Println("macheart not start!")
			return
		}

		// 创建守护进程
		if os.Getppid() != 1 {
			err := global.ForkOneself(os.Args)
			if err != nil {
				fmt.Println("macheart starts has failed!")
			} else {
				fmt.Println("macheart has started!")
			}
			return
		}

	}

	// 启动RPC服务
	rpcServer := rpcserver.New(global.Configure["RpcNetwork"], ":"+global.Configure["RpcPort"])
	err := rpcServer.Start()
	if err != nil {
		global.Logger.Fatal("the RPC server starts has failed: %s", err.Error())
		return
	}

	startHeart()
}

//
func startHeart() {

	// set GOMAXPROCS
	cpus, e := strconv.Atoi(global.Configure["CPUs"])
	if e != nil {
		cpus = 1
	}
	runtime.GOMAXPROCS(cpus)

	// 创建监听接口
	serverAddr, err := net.ResolveTCPAddr(global.Configure["ServerNetwork"], global.Configure["ServerListenAddress"])
	checkFatal("macheart resolve service address failed: %s", err)
	serverListener, err := net.ListenTCP(global.Configure["ServerNetwork"], serverAddr)
	checkFatal("macheart create listener failed: %s", err)
	defer serverListener.Close()

	global.Logger.Info("macheart starts successfully!")
	for {
		// 接受客户端链接
		tcpConn, err := serverListener.AcceptTCP()
		checkFatal("macheart accept the client link failed: ", err)

		// 处理新链接
		go func(conn *net.TCPConn) {

			// 客户端用户验证
			ok, userID := verifyUser(tcpConn)
			if !ok {
				tcpConn.Close()
				return
			}
			global.Logger.Info("user %s(%v) is logged in", userID, tcpConn.RemoteAddr())

			// TODO

		}(tcpConn)

	}
}

func stopHeart() {

}

// 验证用户
func verifyUser(conn *net.TCPConn) (result bool, userID string) {
	return true, "TestUser" //TODO
}

// 检查致命错误并中端程序执行
func checkFatal(logf string, err error) {
	if err != nil {
		global.Logger.Fatal(logf, err.Error())
		os.Exit(1)
	}
}
