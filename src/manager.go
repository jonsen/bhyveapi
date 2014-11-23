// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

type vmInfo struct {
	Name     string
	Cpu      int
	Memory   int
	DiskSize int
	Network  string
	NetCard  int
	Ipaddr   string
	Gateway  string
	Os       string
}

func vmManagerCreate(req *http.Request, r render.Render) {
	decoder := json.NewDecoder(req.Body)
	var (
		info vmInfo
		msg  string = "create vm success"
		bv   *Bhyve
	)

	err := decoder.Decode(&info)
	if err != nil {
		msg = err.Error()
		goto mend
	}

	// TODO check input data
	if _, ok := bhyves[info.Name]; ok {
		msg = "vm " + info.Name + " exists"
		goto mend
	}

	fmt.Println(info)

	//
	bv = new(Bhyve)
	bv.Name = info.Name
	bv.Cpu = info.Cpu
	bv.Memory = info.Memory
	bv.DiskSize = info.DiskSize
	bv.Network = info.Network
	bv.NetCard = info.NetCard
	bv.Ipaddr = info.Ipaddr
	bv.Gateway = info.Gateway
	bv.Os = info.Os

	bhyves[bv.Name] = bv

	go bv.InitDisk()

	// save to file
	err = bhyveDataSave()
	if err != nil {
		msg = err.Error()
	}

mend:
	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)

}

func vmManagerUpdate(req *http.Request, r render.Render, p martini.Params) {

}

func vmManagerDestroy(r render.Render, p martini.Params) {
	var msg string

	vm := p["id"]
	bv, ok := bhyves[vm]
	if ok {
		err := bv.Destroy()
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = "vm not exists"
	}

	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)

}

func vmManagerList(r render.Render) {

}

func vmManagerInfo(r render.Render, p martini.Params) {
	vm := p["id"]
	bv, ok := bhyves[vm]
	if ok {
		r.JSON(200, bv)
		return
	}

	ret := &Message{Message: "vm not exists", Code: 200, Url: ""}
	r.JSON(200, ret)

}
