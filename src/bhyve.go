// Copyright 2014, Jonsen Yang.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	//"os"
	"io"
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
	Disks    []string
	Console  string
}

type BhyveStat struct {
	Pid    int
	Status int

	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

var (
	bhyves      = make(map[string]*Bhyve)
	bhyveStatus = make(map[string]*BhyveStat)
)

func bhyveDataLoad() (err error) {

	body, err := ioutil.ReadFile(gcfg.Global.Datafile)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &bhyves)
	if err != nil {
		return
	}

	stBody, err := ioutil.ReadFile(gcfg.Global.Statfile)
	if err != nil {
		return
	}

	err = json.Unmarshal(stBody, &bhyveStatus)
	return
}

func bhyveDataSave() (err error) {
	body, err := json.MarshalIndent(bhyves, "", "    ")
	if err == nil {
		err = ioutil.WriteFile(gcfg.Global.Datafile, body, 0755)
	}

	stBody, err := json.MarshalIndent(bhyveStatus, "", "    ")
	if err == nil {
		err = ioutil.WriteFile(gcfg.Global.Statfile, stBody, 0755)
	}
	return
}

func GetBhyve(vm string) (b *Bhyve, err error) {
	b, ok := bhyves[vm]
	if !ok {
		return nil, errors.New("vm not exists")
	}

	if _, ok := bhyveStatus[vm]; !ok {
		bhyveStatus[vm] = &BhyveStat{}
	}

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

func (b *Bhyve) IsInstalled() (ok bool, err error) {
	// DOS/MBR boot sector
	// : Unix Fast File sys
	out, err := exec.Command("/usr/bin/file", "-s", gcfg.Global.Vmdir+b.Name+"_0.img").Output()
	if err != nil {
		fmt.Println("IsInstalled", err)
		return false, err
	}

	boots := strings.Index(string(out), "boot sector")
	unixfs := strings.Index(string(out), "Unix Fast File")
	if boots > 0 || unixfs > 0 {
		fmt.Println("aleady install os")
		return true, nil
	}

	return false, errors.New("Need Install OS")
}

func (b *Bhyve) Load() (err error) {
	ok, err := b.IsInstalled()
	if !ok {
		return
	}

	size := fmt.Sprintf("%dM", b.Memory)
	err = exec.Command(gcfg.Bhyve.Bhyveload, "-m", size, "-d", gcfg.Global.Vmdir+b.Name+"_0.img",
		//	b.Name).Run()
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
	st := b.GetStat()
	if st.Pid == 0 {
		return errors.New("vm not running")
	}

	err = syscall.Kill(st.Pid, syscall.SIGTERM)
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
		"-AI", "-H", "-P",
		"-s", "0:0," + b.Network,
		"-s", "1:0,lpc",
		"-s", "2:0,virtio-net,tap0",
		"-s", "3:0,virtio-blk," + gcfg.Global.Vmdir + b.Name + "_0.img",
		//"-l", "com1,stdio",
		"-l", "com1,/dev/nmdm_" + b.Name + "_A",
		b.Name,
	}

	fmt.Println(gcfg.Bhyve.Bhyve, args)

	st := b.GetStat()

	go func() {
		fmt.Println("running bhyve")
		cmd := exec.Command(gcfg.Bhyve.Bhyve, args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}
		st.Stdout = stdout

		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
		st.Stdin = stdin
		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println(err)
		}
		st.Stderr = stderr

		err = cmd.Start()
		if err != nil {
			fmt.Println("start run", err)
		} else {
			st.Pid = cmd.Process.Pid
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
	st := b.GetStat()
	if st.Pid == 0 {
		return errors.New("vm not running")
	}

	err = syscall.Kill(st.Pid, syscall.SIGKILL)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (b *Bhyve) Install() (err error) {

	return
}

func (b *Bhyve) GetStat() (s *BhyveStat) {
	s, ok := bhyveStatus[b.Name]
	if !ok {
		s = &BhyveStat{}
		bhyveStatus[b.Name] = s
	}

	return
}
