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

package server

import (
    "net"
    "net/rpc"
    "macheart/global"
)

type RpcServer struct {
    net string
    addr string
    listener *net.TCPListener
}

func New(net, addr string)(*RpcServer){

    rpcServer := new(RpcServer)
    rpcServer.net = net
    rpcServer.addr = addr
    return rpcServer

}

func (rs *RpcServer)Start()error{

    serverAddr, err := net.ResolveTCPAddr(rs.net, rs.addr)
    if err != nil {return err}
    serverListener, err := net.ListenTCP(rs.net, serverAddr)
    if err != nil {return err}

    rs.listener = serverListener

    go func(listener *net.TCPListener) {
        for {
            conn, err := listener.AcceptTCP()
            if err != nil {
                global.Logger.Warn("the RPC server accept the client link failure: %s", err.Error())
                rs.Stop()
                global.Logger.Warn("the RPC server has stopped!")
                return
            }
            rpc.ServeConn(conn)
        }
    }(rs.listener)

    return nil;
}

func (rs *RpcServer)Stop() error{
    return rs.listener.Close()
}
