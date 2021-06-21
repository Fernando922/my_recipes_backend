package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //driver de conexão com o MySql nao é importado automaticamente
)

//Conectar abre a conexão com o banco de dados
func Connect() (*sql.DB, error) {
	connectionString := "fernando:fernando@/my_recipes?charset=utf8&parseTime=True&loc=Local"

	db, error := sql.Open("mysql", connectionString)

	if error != nil {
		return nil, error
	}

	if error = db.Ping(); error != nil {
		return nil, error
	}

	return db, nil
}
