package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" //solo para sqlite
)

// func (u *usuario) datosUsuarios() usuario {
// 	return &u
// }

//BDAgregarMensaje recibe 2 parametros
func BDAgregarMensaje(m Mensaje) (int, error) {
	db, err := sql.Open("sqlite3", "./iasc.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlAdditem := `
	INSERT OR REPLACE INTO mensajes_sinprocesar(
		id,
		destino,
		tipo,
		payload,
		fecha
	) values(?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sqlAdditem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(m.ID, m.Destino, m.Tipo, m.Payload)
	if err2 != nil {
		fmt.Println("Error al agregar mensaje en BD")
		panic(err2)
	}
	return 1, nil

}

//BDAgregarUsuarioConectado agrega un usuario cuando se inicia sesion
func BDAgregarUsuarioConectado(email string) (int, error) {
	db, err := sql.Open("sqlite3", "./iasc.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlAdditem := `INSERT OR REPLACE INTO usuarios_conectados(email) values(?)`

	stmt, err := db.Prepare(sqlAdditem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(email)
	if err2 != nil {
		fmt.Println("Error al agregar usuarios conectado en BD")
		panic(err2)
	}
	return 1, nil

}

//BDEliminarUsuarioConectado elimina a los usuarios conectados
func BDEliminarUsuarioConectado(email string) (int, error) {
	db, err := sql.Open("sqlite3", "./iasc.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlAdditem := `delete from usuarios_conectados where email = ?`

	stmt, err := db.Prepare(sqlAdditem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(email)
	if err2 != nil {
		fmt.Println("Error al eliminar usuarios conectado en BD")
		panic(err2)
	}
	return 1, nil

}

//ListarUsuarios lista usuarios
func ListarUsuarios() []Usuario {

	db, err := sql.Open("sqlite3", "./iasc.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, email, nombre, apellido from usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var usuarios []Usuario

	for rows.Next() {
		var u Usuario

		err = rows.Scan(&u.ID, &u.Email, &u.Nombre, &u.Apellido)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(u)
		usuarios = append(usuarios, u)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return usuarios
}

//Leer la base de datos
func Leer() {

	db, err := sql.Open("sqlite3", "./bd/iasc.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, email, nombre, apellido from usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var u Usuario

		err = rows.Scan(&u.ID, &u.Email, &u.Nombre, &u.Apellido)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(u)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// sqlStmt := `
	// create table foo (id integer not null primary key, name text);
	// delete from foo;
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// tx.Commit()

}
