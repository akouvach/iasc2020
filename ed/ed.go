package ed

//Usuario es la estructura de datos comun
type Usuario struct {
	Nombre   string
	Apellido string
	Email    string
	ID       int
}

//Mensaje es lo que deben pasar los websocket
type Mensaje struct {
	Channel string
	ID      int
	Payload string
}
