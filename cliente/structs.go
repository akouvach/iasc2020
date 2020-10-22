package main

import "github.com/gorilla/websocket"

//Usuario es la estructura de datos comun
type Usuario struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
}

//Suscripcion corresponde a los diferentes suscriptores
type Suscripcion struct {
	Nombre        string
	Participantes []Cliente
}

//Cliente corresponde a los clientes que se conectan al broker
type Cliente struct {
	ID    string
	email string
	ws    *websocket.Conn
}

//Servicio corresponde a los clientes que se conectan al broker
type Servicio struct {
	nombre string
}

//Mensaje es lo que deben pasar los websocket
type Mensaje struct {
	ID       string `json:"id"`
	Destino  string `json:"destino"`
	Tipo     string `json:"tipo"`
	Persiste bool   `json:"persiste"`
	Payload  string `json:"payload"`
}
