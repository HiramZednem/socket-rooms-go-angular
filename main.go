package main

import (
	"log"
	"net/http"

	"github.com/HiramZednem/SOCKET-ROOMS-GO-ANGULAR/ws"
	"github.com/gorilla/mux"
)

var manager *ws.Manager

func main() {
	setEnviroment()
	r := mux.NewRouter()


	r.HandleFunc("/hello", manager.GetHelloWorld).Methods("GET")
	r.HandleFunc("/ws/createRoom", manager.CreateRoom).Methods("POST")
	r.HandleFunc("/ws/getRooms", manager.GetRooms).Methods("GET")

	//Aqui esta la magia del ws
	r.HandleFunc("/ws/joinRoom/{roomId}", manager.JoinRoom).Methods("GET")

	log.Println("The server is running on port 8080")
	log.Fatal( http.ListenAndServe(":8080", r) )	
}

func setEnviroment() {
	hub := ws.NewHub()
	manager = ws.NewManager(hub)
	go hub.Run()
	
}