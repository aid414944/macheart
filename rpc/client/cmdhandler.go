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
    "fmt"
    "errors"
    "net/rpc"
    "macheart/global"
)

type CmdHandler struct {
    rpcClient *rpc.Client
}

func NewCmdHandler(net, addr string)(*CmdHandler, error) {

    rpcClient, err := rpc.Dial(net, addr)
    if err != nil {return nil, err}

    cmdHandler := new(CmdHandler)
    cmdHandler.rpcClient = rpcClient
    return cmdHandler, nil

}

func (ch *CmdHandler)Close()error{
    return ch.rpcClient.Close()
}

func (ch *CmdHandler)Exec(cmd []string)error{

    var ok bool

    switch cmd[1]{

    case "start":
        fmt.Println("had an instance running!")

    case "restart":

        cmd[1] = "stop"
        fmt.Println("stopping macheart...")
        err := ch.rpcClient.Call("CmdHandleServer.Exec", &cmd, &ok) //若停止成功，这里必然会返回错误,但是返回错误不一定表示服务停止成功;这里姑且认为返回错误就表示停止成功
        if err != nil {
            fmt.Println("macheart has stopped!")
        }else{
            fmt.Println("macheart stops has failed!")
            return errors.New(cmd[1] + ": exec failure")
        }

        cmd[1] = "start"
        fmt.Println("starting macheart...")
        err = global.ForkOneself(cmd)
        if err != nil {
            fmt.Println("macheart starts has failed!")
            return err
        }else {
            fmt.Println("macheart has started!")
        }

    case "stop":
        fmt.Println("stopping macheart...")
        err := ch.rpcClient.Call("CmdHandleServer.Exec", &cmd, &ok)
        if err != nil {
            fmt.Println("macheart has stopped!")
        }else{
            fmt.Println("macheart stops has failed!")
            return errors.New(cmd[1] + ": exec failure")
        }

    default:
        fmt.Println("invalid arguments: ", cmd[1])
        return errors.New("invalid arguments: " + cmd[1])

    }

    return nil
}

