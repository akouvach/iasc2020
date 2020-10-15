package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"bytes"
	"encoding/json"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

//Mensaje es lo que deben pasar los websocket
type Mensaje struct {
	ID       string `json:"id"`
	Destino  string `json:"destino"`
	Tipo     string `json:"tipo"`
	Persiste bool   `json:"persiste"`
	Payload  string `json:"payload"`
}

/*
"fmt"

 "bytes"
  "encoding/json"


type MyStruct struct {
	Name string `json:"name"`
  }
  func main() {
	testStruct := MyStruct{"hello world"}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(testStruct)

	fmt.Println(reqBodyBytes.Bytes()) // this is the []byte
	fmt.Println(string(reqBodyBytes.Bytes())) // converted back to show it's your original object
  }
*/
// br "../broker"
func mensajesSinLeer() {
	fmt.Println("Mensajes sin leer")
}
func mostrarMenu(email string) {
	fmt.Println("Message System - " + email)
	fmt.Println("---------------------")
	fmt.Println("Para enviar mensaje: enviar usuario mensaje")
	fmt.Println("Para listar usuarios: listar usuarios regexp")
	fmt.Println("Para leer los mensajes: leer origen")
	fmt.Println("---------------------")
	fmt.Println("Mensajes sin leer")
	fmt.Println("---------------------")
	mensajesSinLeer()

}

func reader(conn *websocket.Conn, wg *sync.WaitGroup) {
	//defer close(done)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		wg.Done()

		reader := strings.NewReader(string(msg))

		dec := json.NewDecoder(reader)
		var m Mensaje
		err = dec.Decode(&m)
		if err != nil {
			log.Fatal("Error en el unMarshall del mensaje recibido en el Cliente")
		}

		log.Printf("\t\t\t r: %s de %s", m.Payload, m.ID)
	}
}

func enviarMensaje(ws *websocket.Conn, m Mensaje) error {
	myBuffer := new(bytes.Buffer)
	json.NewEncoder(myBuffer).Encode(m)
	//fmt.Println("estoy por mandar ", myBuffer.Bytes())
	err := ws.WriteMessage(websocket.TextMessage, myBuffer.Bytes())
	if err != nil {
		return err
	}

	return nil

}

// func listar(ws *websocket.Conn, email string, comandos []string) {
// 	// usuarios := br.ListarUsuarios()

// 	msg := Mensaje{
// 		ID:       email,
// 		Destino:  "usuarios",
// 		Tipo:     "directo",
// 		Persiste: false,
// 		Payload:  comandos[1],
// 	}

// 	err := enviarMensaje(ws, msg)
// 	if err != nil {
// 		fmt.Println("Error al enviar mensaje: ", msg)
// 	}
// 	// El tipo del mensaje podria ser  Directo / grupo / Tema

// }

func leer() {
	fmt.Println("Leer ")

}

func enviar() {

	fmt.Println("Enviar ")
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

		err := enviarMensaje(ws, msg)
		if err != nil {
			fmt.Println("Error al enviar mensaje: ", msg)
		}

		//Espera simulada
		// espera := rand.Intn(3)
		// time.Sleep(time.Duration(espera) * time.Second)
	}

}

func crearYEnviar(i int, cantMensajesxConexion int, wgg *sync.WaitGroup, persiste bool) {

	defer wgg.Done()

	var wg sync.WaitGroup

	email := "akouvach@yahoo.com"
	c := connect(email)
	defer c.Close()

	wg.Add(cantMensajesxConexion)
	//Envio mensajes
	go enviarMensajes(email, c, cantMensajesxConexion, persiste)

	// me quedo esperando los mensajes de vuelta del servidor
	go reader(c, &wg)

	wg.Wait()
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

func main() {

	// cant, err := ClienteAutomatico(2)
	// if err != nil {
	// 	fmt.Println("Errores", err, cant)
	// }

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("debe ingresar el email para identificarse")
	}
	email := args[0]
	fmt.Println("Bienvenido ", email)

	reader := bufio.NewReader(os.Stdin)

	//creo una conexion para este cliente
	//c := connect()

	for {
		mostrarMenu(email)
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
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
			//listar(c, email, cmd)

		case "leer":
			fmt.Println("leer los mensajes")
		case "enviar":
			fmt.Println("enviar mensaje")

		default:
			fmt.Println("Comando no reconocido")

		}

		// 	// for _, com := range comandos {
		// 	// 	fmt.Println(com)

		// 	// }
		// 	// fmt.Println(strings.Compare("hi", text))

		// 	// if strings.Compare("hi", text) == 0 {
		// 	// 	fmt.Println("hello, Yourself")
		// 	// }

	}

	// broker.Enviar("Cola1", "Mensaje1")
	// broker.Leer()

}
