package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/redis.v3"
	//"github.com/go-redis/redis"
)

//UC es la estructura para la lista de usuarios conectados
const UC = "usuariosconectados"

//US es la estructura que mantiene la lista de usuarios
const US = "usuarios"

//SUS es la estructura que mantiene la lista de suscripciones
const SUS = "suscripciones"

//MESSAGES es la estructura que mantiene la lista de mensajes
const MESSAGES = "mensajes"

func connectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

// func (u *usuario) datosUsuarios() usuario {
// 	return &u
// }

//BDAgregarMensaje 1 parametros
func BDAgregarMensaje(m Mensaje) (int, error) {
	rdb := connectRedis()
	defer rdb.Close()

	//pong, err := rdb.Ping().Result()
	fmt.Println("BDAgregarMensaje")

	// 	err = client.Set("name", "Elliot", 0).Err()
	// // if there has been an error setting the value
	// // handle the error
	// if err != nil {
	//     fmt.Println(err)
	// }

	// val, err := client.Get("name").Result()
	// if err != nil {
	//     fmt.Println(err)
	// }

	// fmt.Println(val)

	// json, err := json.Marshal(Author{Name: "Elliot", Age: 25})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = client.Set("id1234", json, 0).Err()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// val, err := client.Get("id1234").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(val)

	return 1, nil

}

//BDAgregarClienteConectado agrega un usuario cuando se inicia sesion
func BDAgregarClienteConectado(email string) (int, error) {

	rdb := connectRedis()
	defer rdb.Close()

	err := rdb.LPush(UC, email).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Agregar usuarios conectado ", email)

	return 1, nil

}

func bdAgregarSuscripcion(canal string, c Cliente) error {
	rdb := connectRedis()
	defer rdb.Close()

	cant, err := rdb.LLen(SUS).Result()
	if err != nil {
		return err
	}

	sus, err2 := rdb.LRange(SUS, 0, cant).Result()
	if err2 != nil {
		return err2
	}

	fmt.Println(sus)

	lExisteCanal := false

	for i, s := range sus {

		var suscrip Suscripcion

		err = json.Unmarshal([]byte(s), &suscrip)
		if err != nil {
			fmt.Println("Error! al unmarshall de Suscripcion")
			return err
		}

		fmt.Println("comparando ", suscrip.Canal, canal)
		if suscrip.Canal == canal {
			//encontre al canal indicado
			fmt.Println("encontre al canal indicado")
			lExisteCanal = true

			//ahora me fijo si esta suscripto
			lEncontrado := false
			for _, p := range suscrip.Participantes {
				if p == c.email {
					//esta en esta suscripcion.
					lEncontrado = true
					break
				}
			}
			if !lEncontrado {
				//Agrego este usuario a la lista
				suscrip.Participantes = append(suscrip.Participantes, c.email)

				//Actualizo la suscripcion en redis
				susc, err := json.Marshal(suscrip)
				if err != nil {
					return err
				}

				rta, err2 := rdb.LSet(SUS, int64(i), string(susc)).Result()
				if err2 != nil {
					return err2
				}
				fmt.Println("Actualizacion de suscripcion:", rta)

			}

		}

	}

	if !lExisteCanal {

		fmt.Println("No existe el canal")
		//tengo que agregar el canal con este participante
		var sus Suscripcion

		sus.Creador = c.email
		sus.ID = int(cant) + 1
		sus.Canal = canal
		sus.Participantes = append(sus.Participantes, c.email)

		suscrip, err := json.Marshal(sus)
		if err != nil {
			return err
		}

		//Actualizo la suscripcion en redis
		rta, err := rdb.LPush(SUS, string(suscrip)).Result()

		if err != nil {
			return err
		}
		fmt.Println("Actualizacion de suscripcion:", rta)

	}

	return nil

}

//BDSuscribir agrega un usuario a un canal
func BDSuscribir(canal string, ch chan string, c Cliente) error {

	rdb := connectRedis()
	defer rdb.Close()

	//Lo agrego a la lista de suscripciones
	err2 := bdAgregarSuscripcion(canal, c)
	if err2 != nil {
		fmt.Println("error e agregar suscripcion")
		return err2
	}
	fmt.Println("se suscribio al canal en redis ", canal, " correctamente")

	//Me suscribo a un canal con mi mail
	pubsub, err := rdb.Subscribe(canal)
	if err != nil {
		fmt.Println("error al suscribirse al canal")
	}

	fmt.Println("Suscrpcion correcta pubsub", pubsub)

	for {
		mess, err := pubsub.ReceiveMessage()
		if err != nil {
			fmt.Println("error al recibir la pubsub", err)
			break
		}
		//envio el mensage por el canal

		ch <- mess.String()
	}

	return nil

}

//BDClientesPorCanal te devuelve la lista de clientes asociados a un canal
func BDClientesPorCanal(canal string, email string) ([]string, error) {

	rdb := connectRedis()
	defer rdb.Close()

	var clie []string

	cant, err := rdb.LLen(SUS).Result()
	if err != nil {
		return clie, err
	}

	clientes, err := rdb.LRange(SUS, 0, cant).Result()
	if err != nil {
		return clie, err
	}

	fmt.Println("suscripciones de ", canal, ": ", clientes)

	for _, c := range clientes {
		if c != email {
			clie = append(clie, c)
		}
	}

	return clie, nil

}

//BDCanalesSuscriptos son las suscripciones existentes
func BDCanalesSuscriptos(email string) ([]Suscripcion, error) {

	rdb := connectRedis()
	defer rdb.Close()

	var usuSus []Suscripcion

	cant, err := rdb.LLen(SUS).Result()
	if err != nil {
		return usuSus, err
	}

	suscripciones, err := rdb.LRange(SUS, 0, cant).Result()
	if err != nil {
		return usuSus, err
	}

	fmt.Println("suscripciones", suscripciones)

	for _, s := range suscripciones {

		var auxSus Suscripcion

		err = json.Unmarshal([]byte(s), &auxSus)
		if err != nil {
			fmt.Println("Error! al unmarshall de Suscripcion")
			return usuSus, err
		}

		for _, p := range auxSus.Participantes {
			if p == email {
				//esta en esta suscripcion.  Agrego la suscripcion
				usuSus = append(usuSus, auxSus)
				break
			}
		}
	}

	return usuSus, nil

}

//BDMensajesPorCanal son los mensajes de los canales suscriptos
//del usuario
func BDMensajesPorCanal(email string) ([]MensajeCanal, error) {

	rdb := connectRedis()
	defer rdb.Close()

	var mensajes []MensajeCanal

	cant, err := rdb.LLen(MESSAGES).Result()
	if err != nil {
		return mensajes, err
	}

	ms, err2 := rdb.LRange(MESSAGES, 0, cant).Result()
	if err2 != nil {
		return mensajes, err2
	}

	fmt.Println(ms)

	for _, s := range ms {

		var mc MensajeCanal

		err = json.Unmarshal([]byte(s), &mc)
		if err != nil {
			fmt.Println("Error! al unmarshall del mensjae del Canal")
			return mensajes, err
		}

		if mc.ID == email {
			mensajes = append(mensajes, mc)
		}

	}

	return mensajes, nil

}

//BDEliminarUsuarioConectado elimina a los usuarios conectados
func BDEliminarUsuarioConectado(email string) (int, error) {

	rdb := connectRedis()
	defer rdb.Close()

	// pong, err := client.Ping().Result()
	// fmt.Println("BDEliminarUsuarioConectado", pong, err)

	//client.LRange(UC).Err()

	var cont int64 = 0

	val, err := rdb.LRem(UC, cont, email).Result()
	if err != nil {
		return 0, err
	}

	fmt.Println("Se elimino un usuario conectado ", email, val)

	// cant, err := rdb.LLen(UC).Result()
	// if err != nil {
	// 	return 0, err
	// }

	// usuarios, err := rdb.LRange(UC, 0, cant).Result()
	// if err != nil {
	// 	return 0, err
	// }

	// fmt.Println(usuarios)

	// for i, v := range usuarios {
	// 	if v == email {
	// 		// Found!
	// 		var cont int64 = 0
	// 		var a interface{}

	// 	}
	// }

	return 1, nil

}

//BDListarUsuarios lista usuarios
func BDListarUsuarios() ([]Usuario, error) {

	rdb := connectRedis()
	defer rdb.Close()

	var usuarios []Usuario

	// pong, err := client.Ping().Result()

	cant, err := rdb.LLen(US).Result()
	if err != nil {
		return usuarios, err
	}

	users, err := rdb.LRange(US, 0, cant).Result()
	if err != nil {
		return usuarios, err
	}

	for _, v := range users {
		var u Usuario
		err = json.Unmarshal([]byte(v), &u)
		if err != nil {
			return usuarios, err
		}
		usuarios = append(usuarios, u)
	}

	//fmt.Println(users)

	return usuarios, nil
}

//BDAgregarUsuarios sirve para precargar la base
func BDAgregarUsuarios(cant int) error {

	
	rdb := connectRedis()
	defer rdb.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	c, err := rdb.LLen(US).Result()
	if err != nil {
		return err
	}

	for i := 1; i <= cant; i++ {
		var usu Usuario

		usu.Apellido = fmt.Sprintf("Kouvach%d", i)
		usu.Email = fmt.Sprintf("akouvach@yahoo.com%d", i)
		usu.Nombre = fmt.Sprintf("Andres%d", i)
		usu.ID = int(c) + i

		usuario, err := json.Marshal(usu)
		if err != nil {
			return err
		}

		err = rdb.LPush(US, string(usuario)).Err()
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("Usuario agregado ", usuario)

	}

	// pong, err := client.Ping().Result()

	// err = json.Unmarshal([]byte(users), &usuarios)
	// if err != nil {
	// 	return usuarios, err
	// }

	// for i, v := range usuarios {
	// 	if v == email {
	// 		// Found!
	// 		var cont int64 = 0
	// 		var a interface{}

	// 	}
	// }

	return nil
}

//BDLeer la base de datos
func BDLeer() {

	client := connectRedis()
	pong, err := client.Ping().Result()
	fmt.Println("Leer", pong, err)

}

//BDEnviarMensajeCanal envía mensaje a un canal
func BDEnviarMensajeCanal(canal string, mensaje string) {

	client := connectRedis()
	defer client.Close()

	cant := client.Publish(canal, mensaje)

	fmt.Println("mensae publicado en el canal..", cant, canal, mensaje)



}