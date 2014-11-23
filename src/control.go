// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func vmControlStart(r render.Render, p martini.Params) {
	var msg string

	vm := p["id"]
	bv, ok := bhyves[vm]
	if ok {
		err := bv.Start()
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = "vm not exists"
	}

	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)
}

func vmControlReboot(r render.Render, p martini.Params) {
	var msg string

	vm := p["id"]
	bv, ok := bhyves[vm]
	if ok {
		err := bv.Stop()
		err = bv.Start()
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = "vm not exists"
	}

	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)
}

func vmControlHalt(r render.Render, p martini.Params) {
	var msg string

	vm := p["id"]
	bv, ok := bhyves[vm]
	if ok {
		err := bv.Halt()
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = "vm not exists"
	}

	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)
}
