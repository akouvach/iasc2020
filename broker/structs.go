package main

import "github.com/gorilla/websocket"

//Usuario es la estructura de datos comun
type Usuario struct {
	Nombre   string
	Apellido string
	Email    string
	ID       int
}

//Cliente corresponde a los clientes que se conectan al broker
type Cliente struct {
	ID    string
	email string
	ws    *websocket.Conn
}

//Mensaje es lo que deben pasar los websocket
type Mensaje struct {
	ID       string `json:"id"`
	Destino  string `json:"destino"`
	Tipo     string `json:"tipo"`
	Persiste bool   `json:"persiste"`
	Payload  string `json:"payload"`
}
