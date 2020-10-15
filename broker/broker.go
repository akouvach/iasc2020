package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

//  ListarUsuarios es para leer la base de datos
// func ListarUsuarios() []ed.Usuario {
// 	return bd.ListarUsuarios()
// }

//Clientes son el conjunto de cliente conectados
var Clientes = make(map[string]Cliente)
var mClientes sync.Mutex

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

			eliminarCliente(client)
			BDEliminarUsuarioConectado(client.email)
			client.ws.Close()
			fmt.Println("Cliente eliminado.  Restantes:", len(Clientes))

			return
		}
		reader := strings.NewReader(string(msg))

		dec := json.NewDecoder(reader)
		var m Mensaje
		err = dec.Decode(&m)
		if err != nil {
			log.Fatal("Error en el unMarshall del mensaje recibido")
		}

		if m.Persiste {
			fmt.Println("persistiendo")
			_, err := BDAgregarMensaje(m)
			if err != nil {
				fmt.Println("Error al persistir")
			}
		}

		//fmt.Println("Mensaje del Cliente", m)
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

func agregarCliente(c Cliente) {
	mClientes.Lock()

	Clientes[c.ID] = c
	// fmt.Println("-------Clientes----------")
	// for _, cli := range Clientes {
	// 	fmt.Println(cli.ID)
	// }
	// fmt.Println("-------------------------")

	mClientes.Unlock()

}

func eliminarCliente(c Cliente) {
	mClientes.Lock()

	delete(Clientes, c.ID)

	mClientes.Unlock()

}

func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//fmt.Println("nuevo cliente", vars["email"])

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	myUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("No se pudo generar el uuid")
	}
	// lo agrego a mi lista de usuarios conectados
	var nc Cliente
	nc.ID = myUUID.String()
	nc.email = vars["email"]
	nc.ws = ws

	agregarCliente(nc)
	BDAgregarUsuarioConectado(nc.email)

	log.Println("Nuevo Cliente conectado", len(Clientes))
	//fmt.Println(Clientes)

	reader(nc)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homePage)
	r.HandleFunc("/ws/{email}", wsEndPoint)

	//fmt.Println(ListarUsuarios())
	fmt.Println("Web socket server en 8080")
	http.ListenAndServe(":8080", r)

}
