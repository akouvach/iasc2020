package broker

import (
	"../bd"
)

//Enviar es para agregar un mensaje al canal
func Enviar(canal string, mensaje string) {
	bd.Agregar(canal, mensaje)
}

func Leer() {
	bd.Leer()
}
