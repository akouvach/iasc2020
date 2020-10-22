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

//BDAgregarUsuarioConectado agrega un usuario cuando se inicia sesion
func BDAgregarUsuarioConectado(email string) (int, error) {

	rdb := connectRedis()
	defer rdb.Close()

	err := rdb.LPush(UC, email).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Agregar usuarios conectado ", email)

	//usuarios := rdb.LRange(UC, 0, 10)
	// if err != nil {
	// 	if err == redis.Nil {
	// 		fmt.Println("key does not exists")
	// 		return 0, err
	// 	}
	// 	return 0, err
	// }
	//fmt.Println("usuarios", usuarios)

	return 1, nil

}

//BDSuscribir agrega un usuario cuando se inicia sesion
func BDSuscribir(canal string, c chan string) error {

	rdb := connectRedis()
	defer rdb.Close()

	//Me suscribo a un canal con mi mail
	pubsub, err := rdb.Subscribe(canal)
	if err != nil {
		fmt.Println("error al suscribirse al canal")
	}

	fmt.Println("Suscrpcion correcta", pubsub)

	for {
		mess, err := pubsub.ReceiveMessage()
		if err != nil {
			fmt.Println("error al recibir la pubsub", err)
			break
		}
		//envio el mensage por el canal
		c <- mess.String()
	}

	return nil

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

	fmt.Println(suscripciones)

	for _, s := range suscripciones {

		var auxSus Suscripcion

		err = json.Unmarshal([]byte(s), &auxSus)
		if err != nil {
			fmt.Println("Error! al unmarshall de Suscripcion")
			return usuSus, err
		}

		for _, p := range auxSus.Participantes {
			if p.email == email {
				//esta en esta suscripcion.  Agrego la suscripcion
				usuSus = append(usuSus, auxSus)
				break
			}
		}
	}

	return usuSus, nil

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

	for i := 1; i <= cant; i++ {
		var usu Usuario

		usu.Apellido = fmt.Sprintf("Kouvach%d", i)
		usu.Email = fmt.Sprintf("akouvach@yahoo.com%d", i)
		usu.Nombre = fmt.Sprintf("Andres%d", i)
		usu.ID = i

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

//Leer la base de datos
func Leer() {

	client := connectRedis()
	pong, err := client.Ping().Result()
	fmt.Println("Leer", pong, err)

}
