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

package client

import (
	"errors"
	"fmt"
	"macheart/global"
	"net/rpc"
)

type CmdHandler struct {
	rpcClient *rpc.Client
	cmds      map[string]func(args []string) error
}

func NewCmdHandler() *CmdHandler {
	ch := new(CmdHandler)
	ch.cmds = make(map[string]func(args []string) error)
	// start
	ch.cmds["start"] = func(args []string) error {
		fmt.Println("had an instance running!")
		return nil
	}
	// restart
	ch.cmds["restart"] = func(args []string) error {
		var ok bool
		// 停止
		args[1] = "stop"
		fmt.Println("stopping macheart...")
		err := ch.rpcClient.Call("CmdHandleServer.Exec", &args, &ok) //若停止成功，这里必然会返回错误,但是返回错误不一定表示服务停止成功;这里姑且认为返回错误就表示停止成功
		if err != nil {
			fmt.Println("macheart has stopped!")
		} else {
			fmt.Println("macheart stops has failed!")
			return errors.New(args[1] + ": exec failed!")
		}
		// 启动
		args[1] = "start"
		fmt.Println("starting macheart...")
		err = global.ForkOneself(args)
		if err != nil {
			fmt.Println("macheart starts has failed!")
			return err
		} else {
			fmt.Println("macheart has started!")
		}

		return nil
	}
	// stop
	ch.cmds["stop"] = func(args []string) error {
		var ok bool
		fmt.Println("stopping macheart...")
		err := ch.rpcClient.Call("CmdHandleServer.Exec", &args, &ok)
		if err != nil {
			fmt.Println("macheart has stopped!")
		} else {
			fmt.Println("macheart stops has failed!")
			return errors.New(args[1] + ": exec failed!")
		}
		return nil
	}

	return ch
}

func (ch *CmdHandler) LinkServer(net, addr string) error {
	rpcClient, err := rpc.Dial(net, addr)
	if err != nil {
		return err
	}
	ch.rpcClient = rpcClient
	return nil
}

func (ch *CmdHandler) Close() error {
	return ch.rpcClient.Close()
}

func (ch *CmdHandler) Exec(cmd []string) error {
	exec, ok := ch.cmds[cmd[1]]
	if !ok {
		fmt.Println("invalid arguments: ", cmd[1])
		return errors.New("invalid arguments: " + cmd[1])
	}
	return exec(cmd)
}
