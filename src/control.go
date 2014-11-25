// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func vmControlStart(r render.Render, p martini.Params) {
	var msg string = "start success"

	vm := p["id"]

	bv, err := GetBhyve(vm)
	if err == nil {
		err = bv.Start()
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = err.Error()
	}

	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)
}

func vmControlReboot(r render.Render, p martini.Params) {
	var msg string = "reboot success"

	vm := p["id"]
	bv, err := GetBhyve(vm)
	if err == nil {
		err = bv.Halt()
		err = bv.Start()
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = err.Error()
	}

	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)
}

func vmControlHalt(r render.Render, p martini.Params) {
	var msg string = "halt success"

	vm := p["id"]
	bv, err := GetBhyve(vm)
	if err == nil {
		err = bv.Halt()
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = err.Error()
	}
	ret := &Message{Message: msg, Code: 200, Url: ""}
	r.JSON(200, ret)
}
