// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
)

type Message struct {
	Message string
	Code    int
	Url     string
}

func intToStr(i int) (s string) {
	return strconv.Itoa(i)
}
