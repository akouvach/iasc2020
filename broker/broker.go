package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//Clientes son el conjunto de cliente conectados
var Clientes = make(map[string]Cliente)
var mClientes sync.Mutex

func agregarCliente(c Cliente) {
	mClientes.Lock()
	Clientes[c.email] = c
	BDAgregarClienteConectado(c.email)
	mClientes.Unlock()
	log.Println("Nuevo Cliente conectado", len(Clientes), c.email)
}

func eliminarCliente(c Cliente) {
	mClientes.Lock()
	delete(Clientes, c.email)
	BDEliminarUsuarioConectado(c.email)
	mClientes.Unlock()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Para suscribirse ir a /ws")
}

func enviarMensaje(m Mensaje, c Cliente, dest string) error {
	mensaje, err3 := json.Marshal(m)
	if err3 != nil {
		fmt.Println("error del Marshal de mensaje", err3)
		return err3
	}

	//busco el cliente conectado a quien eviar el mensaje

	fmt.Println("buscando cliente ", dest)

	v, existe := Clientes[dest]
	if existe {

		fmt.Println("existe")
		//le mando el mensaje
		err4 := v.ws.WriteMessage(websocket.TextMessage, mensaje)
		if err4 != nil {
			// log.Fatal(err)
			fmt.Println("Error enviando mensaje", err4)
			return err4
		}
		fmt.Println("Mensaje enviado...", mensaje)

	} else {
		fmt.Println("no existe")
	}

	return nil
}
func listarUsuarios(m Mensaje, c Cliente) error {
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

	err3 := enviarMensaje(msg, c, "")
	if err3 != nil {
		fmt.Println("Error al enviar menajes", err3)
		return err3
	}

	return nil

}

func obtenerMensajesPorCanal(m Mensaje, c Cliente) error {
	rta, err := BDMensajesPorCanal(m.ID)
	if err != nil {
		fmt.Println("error al traer los mensajes", err)
		return err
	}

	ms, err2 := json.Marshal(rta)
	if err2 != nil {
		fmt.Println("error del Marshal err2", err2)
		return err2
	}

	msg := Mensaje{
		ID:       "server",
		Destino:  m.Payload,
		Tipo:     "auto",
		Persiste: true,
		Payload:  string(ms),
	}

	mensaje, err3 := json.Marshal(msg)
	if err3 != nil {
		fmt.Println("error del Marshal de mensaje", err3)
		return err3
	}

	//fmt.Println("usuarios", mensaje)
	//fmt.Println("usuarios marshall", usuarios)

	err4 := c.ws.WriteMessage(websocket.TextMessage, mensaje)
	if err4 != nil {
		// log.Fatal(err)
		fmt.Println("Error enviando mensaje", err4)
		return err4
	}

	return nil

}

func procesarMensaje(m Mensaje, c Cliente) {

	fmt.Println("procesando mensaje...", m, c)

	// name := [][]byte{msg, []byte("from server")}
	// sep := []byte("-")

	if m.Destino == "" && m.Tipo == "auto" {
		//Se trata de unas funciones automaticas
		//debo leer el payload
		fmt.Println("Procesando", m.Payload)
		switch m.Payload {
		case "listarUsuarios":
			fmt.Println("listar.")
			listarUsuarios(m, c)
		case "obtenerMensajesPorCanal":
			fmt.Println("obteniendo mensajes por canal.")
			obtenerMensajesPorCanal(m, c)

		default:
			fmt.Println("Payload no reconocido")
		}

	} else {
		// Me est[an pidiendo algo por canal indicado en el campo Destino
		enviarMensajeAlCanal(m.Destino, m.Payload, c)
	}

}

func enviarMensajeAlCanal(canal string, msg string, c Cliente) {

	fmt.Println("enviando mensaje al canal..")
	BDEnviarMensajeCanal(canal, msg)
	// clientes, err := BDClientesPorCanal(canal, c.email)
	// if err != nil {
	// 	fmt.Println("Error al traer los suscriptores", err)
	// }
	// for _, clie := range clientes {
	// 	//Le debo enviar el mensaje
	// 	m := Mensaje{
	// 		ID:       c.email,
	// 		Destino:  clie,
	// 		Tipo:     "auto",
	// 		Persiste: true,
	// 		Payload:  msg,
	// 	}

	// 	fmt.Println("Le voy enviar un mensaje a ", clie)
	// 	err := enviarMensaje(m, 0, c, clie)
	// 	if err != nil {
	// 		fmt.Println("Error al envio del mensaje", err)
	// 	}

	// }

}

func suscribirse(canal string, c Cliente) {
	ch := make(chan string)

	go func() {
		err := BDSuscribir(canal, ch, c)
		if err != nil {
			fmt.Println("error al suscribir", err)
		}
	}()

	// Consume messages.
	for msg := range ch {
		fmt.Println("Mensaje recibido", c.email, msg)
		//Ahora se lo tengo que mandar al cliente
		fmt.Println("enviando mensaje al cliente")


		msg := Mensaje{
			ID:       "server",
			Destino:  canal,
			Tipo:     "auto",
			Persiste: true,
			Payload:  msg,
		}
	
		mensaje, err3 := json.Marshal(msg)
		if err3 != nil {
			fmt.Println("error del Marshal de mensaje", err3)
			return
		}
	
		//fmt.Println("usuarios", mensaje)
		//fmt.Println("usuarios marshall", usuarios)
	
		err4 := c.ws.WriteMessage(websocket.TextMessage, mensaje)  //no se si esta bien pone 0 en messagetype
		if err4 != nil {
			// log.Fatal(err)
			fmt.Println("Error enviando mensaje", err4)
		}
		fmt.Println("mensaje enviado correctamente al cliente")
	
	}

}

func suscribirseACanales(c Cliente){
	suscripciones, err := BDCanalesSuscriptos(c.email)
	if err != nil {
		fmt.Println("No se obtuvieron las suscripciones", err)
	}

	// si no esta suscripto a ninguno, lo agrego al general
	if len(suscripciones) == 0 {
		fmt.Println("no tiene suscripciones")
		suscribirse("general", c)
	} else {
		//Me suscribo a cada canal de los que anteriormente se habia suscripto
		for _, s := range suscripciones {
			suscribirse(s.Canal, c)
		}
	}

}

func leerMensajesDeCliente(c Cliente) {

	//Cargo las suscripciones de este Cliente

	go suscribirseACanales(c)


	// Me quedo esperando instrucciones individuales
	for {

		messageType, msg, err := c.ws.ReadMessage()
		if err != nil {

			eliminarCliente(c)

			c.ws.Close()
			fmt.Println("Cliente eliminado.  Restantes:", len(Clientes))

			return
		}

		fmt.Println("recibiendo mensaje de ", c.email)

		reader := strings.NewReader(string(msg))

		dec := json.NewDecoder(reader)
		var m Mensaje
		err = dec.Decode(&m)

		if err != nil {
			log.Fatal("Error en el unMarshall del mensaje recibido")
		}

		//Procesar mensaje
		fmt.Println("Procesando mensaje recibido", messageType)
		go procesarMensaje(m, c)

	}
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

	// lo agrego a mi lista de usuarios conectados
	var nc Cliente
	nc.email = vars["email"]
	nc.ws = ws

	//Agregar a mi lista de Clientes conectados
	agregarCliente(nc)

	leerMensajesDeCliente(nc)
}

func main() {

	BDAgregarUsuarios(2)
	BDLeer()

	r := mux.NewRouter()

	r.HandleFunc("/", homePage)
	r.HandleFunc("/ws/{email}", wsEndPoint)

	//fmt.Println(ListarUsuarios())
	fmt.Println("Web socket server en 8080")
	http.ListenAndServe(":8080", r)

}
