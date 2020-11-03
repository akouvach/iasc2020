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

	"github.com/google/uuid"

	"encoding/json"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var canalActual string

var mensajesPorCanal []MensajeCanal
var mMensajesPorCanal sync.Mutex

func agregarMensaje(mc MensajeCanal) {
	mMensajesPorCanal.Lock()
	mensajesPorCanal = append(mensajesPorCanal, mc)
	mMensajesPorCanal.Unlock()
}

// br "../broker"
func obtenerMensajesPorCanal(c *websocket.Conn, email string) {

	msg := Mensaje{
		ID:       email,
		Destino:  "server",
		Tipo:     "auto",
		Persiste: true,
		Payload:  "obtenerMensajesPorCanal",
	}

	//Envio mensajes
	err := enviarMensaje(c, msg)
	if err != nil {
		fmt.Println("Error al enviar mensaje: ", msg)
	}

}

func agregarMensajesCanal(m Mensaje) {

	var mensajesCanal []MensajeCanal

	err2 := json.Unmarshal([]byte(m.Payload), &mensajesCanal)
	if err2 != nil {
		fmt.Println("Error! al unmarshall del payload", err2)
	}

	for _, m := range mensajesCanal {
		fmt.Println("Agregando..", m)
		agregarMensaje(m)
	}

}

func mostrarMensajesPorCanal() {
	var resumen map[string]int
	for _, m := range mensajesPorCanal {
		v, existe := resumen[m.Canal]
		if !existe {
			resumen[m.Canal] = 1
		} else {
			resumen[m.Canal] = v + 1
		}
	}
	fmt.Println("Resumen de canales")
	fmt.Println(resumen)

}
func mostrarMenu(email string) {
	fmt.Println("Message System - " + email)
	fmt.Println("---------------------")
	fmt.Println("Para abrir chat: canal <nombre> del canal")
	fmt.Println("Para listar usuarios: listar")
	mostrarMensajesPorCanal()

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

		//se lo envio al canal que lo procesara
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

func cargarMensajesLocalmente(m Mensaje) {

	var mensajes []MensajeCanal

	err2 := json.Unmarshal([]byte(m.Payload), &mensajes)
	if err2 != nil {
		fmt.Println("Error! al unmarshall del payload", err2)
	}

	for _, m := range mensajes {
		agregarMensaje(m)
	}
}

func mostrarMensajes(m Mensaje) {

	var mensajes []MensajeCanal

	err2 := json.Unmarshal([]byte(m.Payload), &mensajes)
	if err2 != nil {
		fmt.Println("Error! al unmarshall del payload", err2)
	}

	for _, m := range mensajes {
		fmt.Println(m.ID, m.Mensaje)
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

		//ya tengo el mensaje recibido
		//en funcion del tipo es que se lo mando a otro para procesarlo

		if m.Tipo == "auto" && m.ID == "server" {
			//corresponden a peticiones que se le hicieron al servidor
			fmt.Println("recibido...", m.Destino)
			switch m.Destino {
			case "listarUsuarios":
				mostrarUsuarios(m)
			case "obtenerMensajesPorCanal":
				cargarMensajesLocalmente(m)
			default:
				fmt.Println("funcionalidad no itentificada")
			}

		} else if m.Tipo == "direct" {
			//Es un mensaje que va para algun canal
			//El canal figura en el campo Destino
			fmt.Println("Es de tipo direct")

		} else {
			fmt.Println("no se reconoce el tipo de mensaje...", m)
		}

		//mostrarMenu(email)

	}
}

func connect(email string) *websocket.Conn {

	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/" + email}
	log.Printf("Conectando a  %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("No se pudo conectar con el servidor :", err)
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

	fmt.Println("enviando mensaje al canal")
	err2 := ws.WriteMessage(websocket.TextMessage, msg)
	if err2 != nil {
		fmt.Println("error al enviar mensaje desde el cliente", err2)
		return err2
	}

	fmt.Println("mensaje enviado correctamente.")

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

func mostrarMensajesCanal(canalActual string, c *websocket.Conn, email string) {

	msg := Mensaje{
		ID:       email,
		Destino:  canalActual,
		Tipo:     "auto",
		Persiste: true,
		Payload:  "listarMensajes",
	}

	//Envio mensajes
	err := enviarMensaje(c, msg)
	if err != nil {
		fmt.Println("Error al enviar mensaje: ", msg)
	}

}

func saludar(c *websocket.Conn, email string,  saludo string) {
	// Envia un buenos dias al canal actual
	var mc MensajeCanal

	mc.ID = uuid.New().String()
	mc.Canal = canalActual
	mc.Fecha = time.Now().String()
	mc.Mensaje = saludo

	mcAux, err2 := json.Marshal(mc)
	if err2!=nil{
		fmt.Println("Error al marshaling...", err2)
		return
	}

	var m Mensaje
	m.Destino = canalActual
	m.ID = email
	m.Tipo = "auto"
	m.Payload = string(mcAux)

	err := enviarMensaje(c, m)

	if err!=nil {
		fmt.Println("error al enviar mensaje... ", err)
	}
	fmt. Println("mensaje enviado correctamente")


}
func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("debe ingresar el email para identificarse")
	}
	email := args[0]
	fmt.Println("Bienvenido ", email)

	c := connect(email)
	defer c.Close()

	//Abro un canal que recibira mensajes del servidor
	ch := make(chan string)
	defer close(ch) // cierro el canal cuando termine

	obtenerMensajesPorCanal(c, email)

	canalActual = "general"
	saludar(c, email,"hola mundo")

	go recibirDelServidor(c, ch) //a este canal, llegaran los mensajes del servidor

	go procesarMensajeDelServidor(c, ch, email) // el anterior canal, se los mandara a este, para procesar

	mostrarMenu(email)

mainCicle:
	for {
		fmt.Print(canalActual + ">_ ")

		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, "\r", "", -1)

		if canalActual == "" {

			strs := strings.Split(line, " ")

			com := strs[0]

			fmt.Println("Comando-"+com+"-", len(com))
			switch {
			case strings.Compare(com, "listar") == 0:
				fmt.Println("listar usuarios")
				fmt.Println("-------------------------------------")
				obtenerUsuarios(c, email)
				fmt.Println("---Presione una tecla para continuar---")
				fmt.Scanln()

			case strings.Compare(com, "canal") == 0:
				fmt.Println("Canal", strs[1])
				canalActual = strs[1]

				fmt.Println("---Presione una tecla para continuar---")
				fmt.Scanln()

			case strings.Compare(com, "salir") == 0:

				fmt.Println("Hasta la vista...")
				break mainCicle
			default:
				fmt.Println("Comando no reconocido")
				fmt.Println("---Presione una tecla para continuar---")
				fmt.Scanln()

			}

		} else {
			//Estoy en un canal.  Voy a mostrar los mensajes
			mostrarMensajesCanal(canalActual, c, email)
			fmt.Println("---------------------------------------------------->")
		}

	}

	// broker.Enviar("Cola1", "Mensaje1")
	// broker.Leer()

}
