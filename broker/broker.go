package main

import (
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

func listarUsuarios(m Mensaje, c *Cliente, mt int) error {
	rta, err := BDListarUsuarios()
	if err != nil {
		fmt.Println("error al listar usuarios", err)
		return err
	}

	usuarios, err2 := json.Marshal(rta)
	if err2 != nil {
		fmt.Println("error del Marshal err2", err2)
		return err2
	}

	msg := Mensaje{
		ID:       "server",
		Destino:  "listarUsuarios",
		Tipo:     "auto",
		Persiste: true,
		Payload:  string(usuarios),
	}

	mensaje, err3 := json.Marshal(msg)
	if err3 != nil {
		fmt.Println("error del Marshal de mensaje", err3)
		return err3
	}

	//fmt.Println("usuarios", mensaje)
	//fmt.Println("usuarios marshall", usuarios)

	err4 := c.ws.WriteMessage(mt, mensaje)
	if err4 != nil {
		// log.Fatal(err)
		fmt.Println("Error enviando mensaje", err4)
		return err4
	}

	return nil

}

func procesarMensaje(m Mensaje, c *Cliente, mt int) {

	fmt.Println("procesando mensaje", m, c)

	// name := [][]byte{msg, []byte("from server")}
	// sep := []byte("-")

	if m.Destino == "" && m.Tipo == "auto" {
		//Se trata de unas funciones automaticas
		//debo leer el payload
		fmt.Println("Procesando", m.Payload)
		switch m.Payload {
		case "listarUsuarios":
			fmt.Println("listar.")
			listarUsuarios(m, c, mt)

		default:
			fmt.Println("Payload no reconocido")
		}

	}

}

func suscribirse(canal string, c *Cliente) {
	ch := make(chan string)

	go func() {
		err := BDSuscribir("general", ch)
		if err != nil {
			fmt.Println("error al suscribir", err)
		}
	}()

	// Consume messages.
	for msg := range ch {
		fmt.Println("Mensaje recibido", c.email, msg)
	}

}
func leerMensajesDeCliente(client *Cliente) {

	//Cargo las suscripciones de este Cliente

	suscripciones, err := BDCanalesSuscriptos(client.email)
	if err != nil {
		fmt.Println("No se obtuvieron las suscripciones", err)
	}

	// si no esta suscripto a ninguno, lo agrego al general
	if len(suscripciones) == 0 {
		fmt.Println("no tiene suscripciones")
		go suscribirse("general", client)
	}

	//Me suscribo a cada canal
	for _, s := range suscripciones {
		go suscribirse(s.Nombre, client)
	}

	// Me quedo esperando instrucciones individuales
	for {

		messageType, msg, err := client.ws.ReadMessage()
		if err != nil {

			eliminarCliente(*client)

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

		//Procesar mensaje
		fmt.Println("Procesando", messageType)
		go procesarMensaje(m, client, messageType)

		// // en M, tengo el mensaje.  Lo agrego a la lista de mensajes
		// _, err = BDAgregarMensaje(m)
		// if err != nil {
		// 	fmt.Println("Error al agregar mensaje")
		// }

		//fmt.Println("Mensaje del Cliente", m)

		// name := [][]byte{msg, []byte("from server")}
		// sep := []byte("-")

		// err = client.ws.WriteMessage(messageType, bytes.Join(name, sep))
		// if err != nil {
		// 	// log.Fatal(err)
		// 	fmt.Println(err)
		// 	return
		// }
	}
}

func agregarCliente(c Cliente) {
	mClientes.Lock()
	Clientes[c.ID] = c
	BDAgregarUsuarioConectado(c.email)
	mClientes.Unlock()
}

func eliminarCliente(c Cliente) {
	mClientes.Lock()
	delete(Clientes, c.ID)
	BDEliminarUsuarioConectado(c.email)
	mClientes.Unlock()
}

// func agregarServicio(s Servicio) {
// 	mServicios.Lock()
// 	Servicios[s.nombre] = s
// 	mServicios.Unlock()
// }

// func eliminarServicio(s Servicio) {
// 	mServicios.Lock()
// 	delete(Servicios, s.nombre)
// 	mServicios.Unlock()
// }

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

	log.Println("Nuevo Cliente conectado", len(Clientes))
	//fmt.Println(Clientes)

	leerMensajesDeCliente(&nc)
}

func main() {

	BDAgregarUsuarios(5)

	r := mux.NewRouter()

	r.HandleFunc("/", homePage)
	r.HandleFunc("/ws/{email}", wsEndPoint)

	//fmt.Println(ListarUsuarios())
	fmt.Println("Web socket server en 8080")
	http.ListenAndServe(":8080", r)

}
