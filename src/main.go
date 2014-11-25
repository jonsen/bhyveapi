// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

func main() {
	loadConfig()
	bhyveDataLoad()

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Group("/manager", func(r martini.Router) {
		r.Post("/create", vmManagerCreate)
		r.Post("/update/:id", vmManagerUpdate)
		r.Get("/destroy/:id", vmManagerDestroy)
		r.Get("/list", vmManagerList)
		r.Get("/info/:id", vmManagerInfo)
	})

	m.Group("/control", func(r martini.Router) {
		r.Get("/start/:id", vmControlStart)
		r.Get("/reboot/:id", vmControlReboot)
		r.Get("/halt/:id", vmControlHalt)
	})

	//m.Group("/manager", func(r martini.Router) {
	//})

	m.Group("/status", func(r martini.Router) {
		r.Get("", vmStatus)     // get global status
		r.Get("/:id", vmStatus) // get status by vm name
	})

	m.Get("/console/:id", webConsole)

	port := fmt.Sprintf(":%d", gcfg.Global.Port)
	http.ListenAndServe(port, m)
}
