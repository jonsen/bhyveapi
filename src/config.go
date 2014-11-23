// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

// define config struct
type config struct {
	Global globalConfig
	Auth   authConfig
	Bhyve  bhyveConfig
}

// define global config struct
type globalConfig struct {
	Port     int
	Datafile string
	Vmdir    string
}

// define author config struct
type authConfig struct {
	Enabled bool
	Authkey string
	Clients []string
}

// define bhyve config struct
type bhyveConfig struct {
	Bhyve     string
	Bhyvectl  string
	Bhyveload string
	Bhyvegrub string
}

// init global config
var gcfg config

func loadConfig() {
	if _, err := toml.DecodeFile("etc/bhyveapid.conf", &gcfg); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(gcfg)
}
