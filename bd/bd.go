package bd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" //solo para sqlite
)

type usuario struct {
	nombre   string
	apellido string
	email    string
	id       int
}

// func (u *usuario) datosUsuarios() usuario {
// 	return &u
// }

//Agregar recibe 2 parametros
func Agregar(cola string, mensaje string) {
	var msg = fmt.Sprintf("Agregando %s a la cola %s", mensaje, cola)
	fmt.Println(msg)
}

//Leer la base de datos
func Leer() {

	db, err := sql.Open("sqlite3", "./iasc.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var u usuario

		err = rows.Scan(&u)
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
