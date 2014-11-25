package main

import (
	"bufio"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	//"io"
	"net/http"
)

func webConsole(w http.ResponseWriter, r *http.Request, p martini.Params) {
	vm := p["id"]
	bv, err := GetBhyve(vm)
	if err != nil {
		http.Error(w, "Not a vm exists", 400)
		return
	}

	st := bv.GetStat()

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 401)
		fmt.Println("HandshakeError", err)
		return
	} else if err != nil {
		fmt.Println("Upgrade", err)
		return
	}

	defer ws.Close()

	fmt.Println(bv.Name)
	//client := ws.RemoteAddr()
	//fmt.Println(client)

	messageType := 1

	go func() {
		/*
			for {
				fmt.Println("NextWriter1")
				w, err := ws.NextWriter(2)
				if err != nil {
					fmt.Println("NextWriter", err)
					return
				}
				fmt.Println("NextWriter2")
				if _, err := io.Copy(w, bv.Stderr); err != nil {
					fmt.Println("Copy out", err)
					continue
				}
				fmt.Println("NextWriter3")
			}
		*/

		bf := bufio.NewReader(st.Stderr)
		for {
			fmt.Println("ReadSlice")
			line, err := bf.ReadSlice('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("line", string(line))
			if err := ws.WriteMessage(messageType, line); err != nil {
				fmt.Println("WriteMessage22", err)
				return
			}

		}
	}()

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("bye")
			fmt.Println(err)
			return
		}
		if err := ws.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}

		/*
			w, err := ws.NextWriter(messageType)
			if err != nil {
				fmt.Println("NextWriter", err)
				return
			}
			if _, err := io.Copy(w, bv.Stdout); err != nil {
				fmt.Println("Copy out", err)
				continue
			}
		*/

		/*
			w, err := ws.NextWriter(messageType)
			if err != nil {
				fmt.Println("NextWriter", err)
				return
			}

			_, r, err := ws.NextReader()
			if err != nil {
				fmt.Println("NextReader", err)
				return
			}

			if _, err := io.Copy(w, r); err != nil {
				fmt.Println("Copy echo", err)
				continue
			}
			if _, err := io.Copy(w, bv.Stdout); err != nil {
				fmt.Println("Copy out", err)
				continue
			}
			if _, err := io.Copy(bv.Stdin, r); err != nil {
				fmt.Println("Copy in", err)
				continue
			}
		*/
		//if err := w.Close(); err != nil {
		//	return err
		//}

		/*
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("bye")
				fmt.Println(err)
				return
			}
			if err := ws.WriteMessage(messageType, p); err != nil {
				fmt.Println(err)
				return
			}
		*/

	}
}
