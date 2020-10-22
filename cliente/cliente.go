package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

// br "../broker"
func mensajesSinLeer() {
	fmt.Println("Mensajes sin leer")
}
func mostrarMenu(email string) {
	fmt.Println("Message System - " + email)
	fmt.Println("---------------------")
	fmt.Println("Para enviar mensaje: enviar usuario mensaje")
	fmt.Println("Para listar usuarios: listar")
	fmt.Println("Para leer los mensajes: leer canal")
	fmt.Println("---------------------")
	fmt.Println("Mensajes sin leer")
	fmt.Println("---------------------")
	mensajesSinLeer()

}

func recibirDelServidor(conn *websocket.Conn, c chan string) {
	//defer close(done)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			close(c) //cierro el canal
			return
		}

		//se lo envio al cana que lo procesara
		c <- string(msg)
	}
}

func mostrarUsuarios(m Mensaje) {

	var usuarios []Usuario

	err2 := json.Unmarshal([]byte(m.Payload), &usuarios)
	if err2 != nil {
		fmt.Println("Error! al unmarshall del payload", err2)
	}

	for _, u := range usuarios {
		fmt.Println(u.Email)
	}
}

func procesarMensajeDelServidor(conn *websocket.Conn, c chan string, email string) {
	//defer close(done)
	for msg := range c {
		fmt.Println("recibiendo mensaje...")

		var m Mensaje

		err := json.Unmarshal([]byte(msg), &m)
		if err != nil {
			fmt.Println("Error! al unmarshall del mensaje recibido", err)
		}

		//ya tengo el mensaje recibid
		//en funcion del tipo es que se lo mando a otro
		//para procesarlo

		if m.Tipo == "auto" && m.ID == "server" {
			//corresponden a peticiones que se le hicieron al servidor
			switch m.Destino {
			case "listarUsuarios":
				mostrarUsuarios(m)
			default:
				fmt.Println("funcionalidad no itentificada")
			}

		}

		mostrarMenu(email)

	}
}

func connect(email string) *websocket.Conn {

	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/" + email}
	log.Printf("Conectando a  %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	//fmt.Println(fmt.Sprintf("-------Nueva conexion %d ------", i))

	return c
}

func enviarMensajes(email string, ws *websocket.Conn, cant int, persiste bool) {

	for i := 1; i <= cant; i++ {
		//Mando mensaje
		//fmt.Printf("Enviando mensaje %d de %d \n", i, id)

		msg := Mensaje{
			ID:       email,
			Destino:  "MensajesAutomaticos",
			Tipo:     "directo",
			Persiste: persiste,
			Payload:  "pl " + fmt.Sprintf("%d", i),
		}

		//Espera simulada
		espera := rand.Intn(3)
		time.Sleep(time.Duration(espera) * time.Second)

		err := enviarMensaje(ws, msg)
		if err != nil {
			fmt.Println("Error al enviar mensaje: ", msg)
		}

	}

}

func crearYEnviar(i int, cantMensajesxConexion int, wgg *sync.WaitGroup, persiste bool) {

	defer wgg.Done()

	email := fmt.Sprintf("akouvach@yahoo.com%d", i)
	c := connect(email)
	defer c.Close()

	//Envio mensajes
	go enviarMensajes(email, c, cantMensajesxConexion, persiste)

	ch := make(chan string)
	defer close(ch)
	// me quedo esperando los mensajes de vuelta del servidor
	go recibirDelServidor(c, ch)

	fmt.Println(fmt.Sprintf("Se termino de procesar los mensajes de %d", i))

}

//ClienteAutomatico va a crear la cantidad de conexiones y va a mandar 5 mensajes al
// servidor y luego va a cerrar
func ClienteAutomatico(cantConexiones int, cantMensajesxConexion int, persiste bool) string {

	var wgg sync.WaitGroup

	for i := 1; i <= cantConexiones; i++ {
		//id := time.Now()
		wgg.Add(1)
		go crearYEnviar(i, cantMensajesxConexion, &wgg, persiste)

	}
	fmt.Println("Esperando a que termine todo----")
	wgg.Wait()

	fmt.Println("termino el ciclo de creacion de conexiones--------------")

	return "OK"

}

func enviarMensaje(ws *websocket.Conn, m Mensaje) error {

	msg, err := json.Marshal(m)
	if err != nil {
		fmt.Println("error en el marshal", err)
		return err
	}

	err2 := ws.WriteMessage(websocket.TextMessage, msg)
	if err2 != nil {
		fmt.Println("error al enviar mensaje desde el cliente", err2)
		return err2
	}

	return nil

}

func obtenerUsuarios(c *websocket.Conn, email string) {

	msg := Mensaje{
		ID:       email,
		Destino:  "",
		Tipo:     "auto",
		Persiste: true,
		Payload:  "listarUsuarios",
	}

	//Envio mensajes
	err := enviarMensaje(c, msg)
	if err != nil {
		fmt.Println("Error al enviar mensaje: ", msg)
	}

}

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("debe ingresar el email para identificarse")
	}
	email := args[0]
	fmt.Println("Bienvenido ", email)

	myReader := bufio.NewReader(os.Stdin)

	c := connect(email)
	defer c.Close()

	//Abro un canal que recibira mensajes del servidor
	ch := make(chan string)
	defer close(ch) // cierro el canal cuando termine

	go recibirDelServidor(c, ch) //a este canal, llegaran los mensajes del servidor

	go procesarMensajeDelServidor(c, ch, email) // el anterior canal, se los mandara a este, para procesar

	for {
		mostrarMenu(email)
		fmt.Print("-> ")
		text, _ := myReader.ReadString('\n')
		// convert CRLF to LF
		//text = strings.Replace(text, "\n", " ", -1)

		var cmd []string
		comandos := strings.Split(text, " ")
		for _, element := range comandos {
			cmd = append(cmd, strings.ToLower(strings.Trim(element, "")))
		}

		fmt.Println(cmd)

		com := string(cmd[0])
		fmt.Println(com, len(com))

		switch com {
		case "listar":
			fmt.Println("listar usuarios")
			fmt.Println("-------------------------------------")
			obtenerUsuarios(c, email)

		case "leer":
			fmt.Println("leer los mensajes")
		case "enviar":
			fmt.Println("enviar mensaje")
		case "salir":

			fmt.Println("Hasta la vista...")
			break
		default:
			fmt.Println("Comando no reconocido")

		}

	}

	// broker.Enviar("Cola1", "Mensaje1")
	// broker.Leer()

}
