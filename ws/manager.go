package ws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // Se le puede meter aca directo las direcciones para permitir ciertos origenes
		return true
	},
}

type Manager struct {
	hub *Hub
}

func NewManager(h *Hub) *Manager {
	return &Manager{
		hub: h,
	}
}

func (h *Manager) GetHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World from manager")
}

// El apartado que me pase el ID hay que buscar una libreria que genere cadenas aleatorias
// como estas d8q1p-b0hmk-wgxf4, para que el front no las este mandando, o meter algo en BD

type CreateRoomReq struct {
	Name string `json:"name"`
}


var ID  int = 0

func (h *Manager) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room CreateRoomReq

	// TODO: Meter restriccion en que se siga el formato a pie de la letra
	reqBody, err := ioutil.ReadAll( r.Body )
	if err != nil {
		fmt.Fprint(w, "Something went wrong")
	}
	
	json.Unmarshal( reqBody, &room )

	h.hub.Rooms[ strconv.Itoa(ID) ] = &Room{
		ID:      strconv.Itoa(ID),
		Name:    room.Name,
		Clients: make(map[string]*Client),
	}
	ID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)

}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Manager) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rooms)
}



// // La magia esta en esta parte, aqui se hace la conexion y se hace la conexion con el hub.
func (h *Manager) JoinRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Print("New connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	vars := mux.Vars(r)
	roomID := vars["roomId"]
	clientID := r.URL.Query().Get("userId")
	username := r.URL.Query().Get("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	// En estas partes nos conectamos a los canales que esta leyendo en ciclo infinito el hub
	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}




// type ClientRes struct {
// 	ID       string `json:"id"`
// 	Username string `json:"username"`
// }

// func (h *Manager) GetClients(c *gin.Context) {
// 	var clients []ClientRes
// 	roomId := c.Param("roomId")

// 	if _, ok := h.hub.Rooms[roomId]; !ok {
// 		clients = make([]ClientRes, 0)
// 		c.JSON(http.StatusOK, clients)
// 	}

// 	for _, c := range h.hub.Rooms[roomId].Clients {
// 		clients = append(clients, ClientRes{
// 			ID:       c.ID,
// 			Username: c.Username,
// 		})
// 	}

// 	c.JSON(http.StatusOK, clients)
// }
