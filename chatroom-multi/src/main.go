// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")
var house = make(map[string]*Hub) //每個房間都有hub

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	// if r.URL.Path != "/" {
	// 	http.Error(w, "Not found", http.StatusNotFound)
	// 	return
	// }
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	// hub := newHub()
	// go hub.run() //不管多少客戶端連接，只開了一個goroutine

	r := mux.NewRouter()               //建立mux實體 router
	r.HandleFunc("/{room}", serveHome) //把HTML文件返回給請求方(瀏覽器)
	r.HandleFunc("/ws/{room}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)       //套件會開一張map 讓你儲存房間號碼
		roomId := vars["room"]    //動態取得URL的room
		room, ok := house[roomId] //如果hub中有那個房間

		var hub *Hub
		if ok {
			hub = room
		} else {
			hub = newHub(roomId)
			house[roomId] = hub
			go hub.run()
		}
		serveWs(hub, w, r)
	})
	// server := &http.Server{
	// 	Addr:              *addr,
	// 	ReadHeaderTimeout: 3 * time.Second,
	// }
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
