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
	"macheart/global"
	"net/rpc"
	"os"
)

type CmdHandleServer struct {
	cmds map[string]func(args []string) error
}

func NewCmdHandleServer() *CmdHandleServer {
	chs := new(CmdHandleServer)
	chs.cmds = make(map[string]func(args []string) error)
	// stop
	chs.cmds["stop"] = func(args []string) error {
		//stopHeart()  //善后处理
		global.Logger.Info("macheart has stopped!")
		os.Exit(0)
		return nil
	}

	return chs
}

func (chs *CmdHandleServer) Exec(cmd *[]string, ok *bool) error {
	*ok = true
	exec, yes := chs.cmds[(*cmd)[1]]
	if !yes {
		*ok = false
		return nil
	}
	return exec(*cmd)
}

func init() {
	cmdHandleServer := NewCmdHandleServer()
	rpc.Register(cmdHandleServer)
}
