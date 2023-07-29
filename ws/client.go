package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message // Los mensajes llegan en el canal de los mensajes, si se recibe un mensaje
		if !ok {
			return
		}

		c.Conn.WriteJSON(message) //Se escribe en la conexion existente
		//Esta onda se le manda al cliente
		/*
			Content
			RoomId
			Username
		*/
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		//Pum, aca el cliente me manda un mensaje, que es un []char
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		//creo un objeto mensaje
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		//Y lo mando a que se mande a los demas usuarios conectados
		hub.Broadcast <- msg
	}
}
