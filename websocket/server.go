// Copyright 2011 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package websocket

import (
	"errors"
	"github.com/garyburd/t2/web"
	"strings"
)

func inListHeader(header web.Header, key string, value string) bool {
	for _, v := range header.GetList(key) {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
}

// Upgrade upgrades the HTTP connection to the WebSocket protocol. 
func Upgrade(resp web.Response, req *web.Request, subProtocol string) (*Conn, error) {

	if req.Method != "GET" {
		return nil, &web.Error{Status: web.StatusMethodNotAllowed}
	}

	if "13" != req.Header.Get(web.HeaderSecWebSocketVersion) {
		return nil, &web.Error{
			Status: web.StatusBadRequest,
			Reason: errors.New("websocket: version != 13")}
	}

	if !inListHeader(req.Header, web.HeaderConnection, "upgrade") {
		return nil, &web.Error{
			Status: web.StatusBadRequest,
			Reason: errors.New("websocket: connection header != upgrade")}
	}

	if !inListHeader(req.Header, web.HeaderUpgrade, "websocket") {
		return nil, &web.Error{
			Status: web.StatusBadRequest,
			Reason: errors.New("websocket: upgrade != websocket")}
	}

	netConn, br, err := resp.Hijack()
	if err != nil {
		return nil, err
	}

	conn, err := NewServerConn(netConn, br, 4096, subProtocol, req.Header.Get(web.HeaderSecWebSocketKey))
	if err != nil {
		netConn.Close()
	}

	return conn, err
}
