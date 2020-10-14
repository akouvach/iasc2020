package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

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
	fmt.Println("Para listar usuarios: listar usuarios regexp")
	fmt.Println("Para leer los mensajes: leer origen")
	fmt.Println("---------------------")
	fmt.Println("Mensajes sin leer")
	fmt.Println("---------------------")
	mensajesSinLeer()

}

func listar(ws *websocket.Conn) {
	// usuarios := br.ListarUsuarios()

	err := ws.WriteMessage(websocket.TextMessage, []byte("hola"))
	if err != nil {
		log.Println("write:", err)
		return
	}

}

func leer() {
	fmt.Println("Leer ")

}

func enviar() {

	fmt.Println("Enviar ")
}

func connect() *websocket.Conn {

	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("Conectando a  %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func enviarMensajes(id time.Time, wg *sync.WaitGroup, ws *websocket.Conn, cant int) {

	defer wg.Done()

	fmt.Printf("Conexion %s empezando\n", id.String())

	for i := 1; i <= cant; i++ {
		//Mando mensaje

		err := ws.WriteMessage(websocket.TextMessage, []byte(id.String()+"-hola"))
		if err != nil {
			log.Println("write:", err)
			return
		}
		//rand.Seed(99)
		espera := rand.Intn(3)
		fmt.Println("esperando %i segundos para mandar mensaje", espera)
		time.Sleep(time.Duration(espera) * time.Second)
	}

	fmt.Printf("Conexion %s done\n", id.String())

}

//ClienteAutomatico va a crear la cantidad de conexiones y va a mandar 5 mensajes al
// servidor y luego va a cerrar
func ClienteAutomatico(cantidad int) (int, error) {

	n := 0
	var wg sync.WaitGroup

	for i := 1; i <= cantidad; i++ {
		id := time.Now()
		fmt.Println("------------Abro una conexion " + id.String() + " ----------")

		c := connect()
		defer c.Close()

		// me quedo esperando los mensajes de vuelta del servidor
		go func() {
			//defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recibiendo en "+id.String()+" : %s", message)
			}
		}()

		wg.Add(1)
		go enviarMensajes(id, &wg, c, 5)
		n++
	}

	wg.Wait()
	fmt.Println("termino el ciclo")
	return n, nil

}

func main() {

	ClienteAutomatico(2)

	// args := os.Args[1:]
	// if len(args) == 0 {
	// 	log.Fatal("debe ingresar el email para identificarse")
	// }
	// fmt.Println(args[0])

	// reader := bufio.NewReader(os.Stdin)
	// //fmt.Println("ingrese su correo")
	// //email, _ := reader.ReadString('\n')
	// email := args[0]
	// mostrarMenu(email)
	// for {
	// 	fmt.Print("-> ")
	// 	text, _ := reader.ReadString('\n')
	// 	// convert CRLF to LF
	// 	//text = strings.Replace(text, "\n", " ", -1)

	// 	comandos := strings.Split(text, " ")

	// 	com := strings.Trim(strings.ToLower(comandos[0]), " ")

	// 	switch {
	// 	case strings.Compare(com, "listar") == 0:
	// 		listar(c)
	// 	case strings.Compare(com, "leer") == 0:
	// 		fmt.Println("leer los mensajes")
	// 	case strings.Compare(com, "enviar") == 0:
	// 		fmt.Println("enviar mensaje")

	// 	default:
	// 		fmt.Println("Comando no reconocido")
	// 		mostrarMenu(email)

	// 	}

	// 	// for _, com := range comandos {
	// 	// 	fmt.Println(com)

	// 	// }
	// 	// fmt.Println(strings.Compare("hi", text))

	// 	// if strings.Compare("hi", text) == 0 {
	// 	// 	fmt.Println("hello, Yourself")
	// 	// }

	// }

	// broker.Enviar("Cola1", "Mensaje1")
	// broker.Leer()

}
