// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	//"os"
	"os/exec"
	"syscall"
)

type Bhyve struct {
	Name     string
	Cpu      int
	Memory   int
	DiskSize int
	Network  string
	NetCard  int
	Ipaddr   string
	Gateway  string
	Os       string
	Id       int
	Pid      int
	Status   int
	Disks    []string
	Console  string
}

var bhyves map[string]*Bhyve

func bhyveDataLoad() (err error) {

	body, err := ioutil.ReadFile(gcfg.Global.Datafile)
	if err != nil {
		bhyves = make(map[string]*Bhyve)
		return
	}

	err = json.Unmarshal(body, &bhyves)
	return
}

func bhyveDataSave() (err error) {
	body, err := json.MarshalIndent(bhyves, "", "    ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(gcfg.Global.Datafile, body, 0755)
	return
}

func (b *Bhyve) InitDisk() (err error) {
	fmt.Println("init disk")
	size := fmt.Sprintf("%dG", b.DiskSize)
	err = exec.Command("/usr/bin/truncate", "-s", size, b.Name+"_0.img").Run()
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (b *Bhyve) InitNetwork() (err error) {

	return
}

func (b *Bhyve) Load() (err error) {
	size := fmt.Sprintf("%dM", b.Memory)
	err = exec.Command(gcfg.Bhyve.Bhyveload, "-m", size, "-d", b.Name+"_0.img",
		"-c", "/dev/nmdm_"+b.Name+"_A", b.Name).Run()
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (b *Bhyve) Destroy() (err error) {
	err = exec.Command(gcfg.Bhyve.Bhyvectl, "--destroy", "--vm="+b.Name).Run()
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (b *Bhyve) Halt() (err error) {
	if b.Pid == 0 {
		return errors.New("vm not running")
	}

	err = syscall.Kill(b.Pid, syscall.SIGTERM)
	if err != nil {
		fmt.Println(err)
		return
	}

	//
	return
}

// /usr/sbin/bhyve -c 2 -m 256 -A -H -P \
// -s 0:0,hostbridge \
// -s 1:0,virtio-net,tap0 \
// -s 2:0,ahci-hd,./vm0.img \
// -s 31,lpc -l com1,stdio \
// vm0
func (b *Bhyve) Start() (err error) {
	/*
		// Parent process*
		rCmdIn, lCmdOut, err := os.Pipe() // Pipe to write from parent to remote
		child's stdin.
		exitOnErr("1. os.Pipe(): ", err)
		lCmdIn, rCmdOut, err := os.Pipe() // Pipe to read from remote cmd's stdout.
		exitOnErr("2. os.Pipe(): ", err)

		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{rCmdIn, rCmdOut, os.Stderr}
		pid, err := os.StartProcess(ssh, args, &procAttr)
	*/

	var (
		args []string
	)

	b.Destroy()
	err = b.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	args = []string{"-c", intToStr(b.Cpu),
		"-m", intToStr(b.Memory),
		"-A", "-H", "-P",
		"-s", "0:0," + b.Network,
		"-s", "1:0,lpc",
		"-s", "2:0,virtio-net,tap0",
		"-s", "3:0,ahci-hd," + gcfg.Global.Vmdir + b.Name + "_0.img",
		"-l", "com1,/dev/nmdm_" + b.Name + "_A",
		b.Name,
	}

	fmt.Println(gcfg.Bhyve.Bhyve, args)

	go func() {
		fmt.Println("running bhyve")
		//var procAttr os.ProcAttr
		//p, err := os.StartProcess(gcfg.Bhyve.Bhyve, args, &procAttr)
		cmd := exec.Command(gcfg.Bhyve.Bhyve, args...)
		err = cmd.Start()
		if err != nil {
			fmt.Println("start run", err)
		} else {
			b.Pid = cmd.Process.Pid
			fmt.Println("bhyve running now...")

			err = cmd.Wait()
			if err != nil {
				fmt.Println("wait error", err)
			}
		}

		fmt.Println("process exit")
	}()
	return
}

func (b *Bhyve) Stop() (err error) {
	if b.Pid == 0 {
		return errors.New("vm not running")
	}

	err = syscall.Kill(b.Pid, syscall.SIGKILL)
	if err != nil {
		fmt.Println(err)
	}
	return
}
