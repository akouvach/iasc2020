package bd

import (
	"database/sql"
	"fmt"
	"log"

	ed "../ed"

	_ "github.com/mattn/go-sqlite3" //solo para sqlite
)

// func (u *usuario) datosUsuarios() usuario {
// 	return &u
// }

//Agregar recibe 2 parametros
func Agregar(cola string, mensaje string) {
	var msg = fmt.Sprintf("Agregando %s a la cola %s", mensaje, cola)
	fmt.Println(msg)
}

//ListarUsuarios lista usuarios
func ListarUsuarios() []ed.Usuario {

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

	var usuarios []ed.Usuario

	for rows.Next() {
		var u ed.Usuario

		err = rows.Scan(&u.Id, &u.Email, &u.Nombre, &u.Apellido)
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
		var u ed.Usuario

		err = rows.Scan(&u.Id, &u.Email, &u.Nombre, &u.Apellido)
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
