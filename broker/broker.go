package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// bd "../bd"
// ed "../ed"
//Enviar es para agregar un mensaje al canal
// func Enviar(canal string, mensaje string) {
// 	bd.Agregar(canal, mensaje)
// }

//Leer es para leer la base de datos
// func Leer() {
// 	bd.Leer()
// }

//ListarUsuarios es para leer la base de datos
// func ListarUsuarios() []ed.Usuario {
// 	return bd.ListarUsuarios()
// }
//Cliente corresponde a los clientes que se conectan al broker
type Cliente struct {
	ID int
	ws *websocket.Conn
}

//Clientes son el conjunto de cliente conectados
var Clientes = make(map[int]Cliente)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Para suscribirse ir a /ws")
}

func reader(client Cliente) {
	for {

		messageType, msg, err := client.ws.ReadMessage()
		if err != nil {
			fmt.Println("Eliminando cliente")
			delete(Clientes, client.ID)
			//log.Fatal(err)
			client.ws.Close()
			fmt.Println("Clientes actuales")
			fmt.Println(Clientes)
			return
		}
		log.Println(string(msg))
		name := [][]byte{msg, []byte("from server")}
		sep := []byte("-")

		err = client.ws.WriteMessage(messageType, bytes.Join(name, sep))
		if err != nil {
			// log.Fatal(err)
			fmt.Println(err)
			return
		}
	}
}
func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Nuevo Cliente conectado")
	// lo agrego a mi lista de usuarios conectados
	var nc Cliente
	nro := len(Clientes)

	nc.ID = nro
	nc.ws = ws

	Clientes[nro] = nc
	fmt.Println(Clientes)

	reader(nc)
}
func setUpRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndPoint)
}
func main() {
	setUpRoutes()
	//fmt.Println(ListarUsuarios())
	fmt.Println("Web socket server en 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
