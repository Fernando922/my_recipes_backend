package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"recipes-backend/database"
)

type recipe struct {
	ID           uint32 `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Prepare_mode string `json:"prepare_mode"`
}

type recipeType struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type ingredient struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	RecipeId string `json:"recipeId"`
}

func CreateRecipeType(rw http.ResponseWriter, r *http.Request) {

	//lendo o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Request body read error"))
		return
	}

	//convertendo o json do body da requisição para struct
	var recipeType recipeType
	if error = json.Unmarshal(requestBody, &recipeType); error != nil {
		rw.Write([]byte("Error to convert type to struct"))
		return
	}
	fmt.Println(recipeType)

	//abrindo conexão com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close() //só será executado no final da função

	//Criando o prepare statement (comando de inserção, para evitar sql injection)
	statement, error := db.Prepare("insert into Types (name, image) values (?,?)")
	if error != nil {
		rw.Write([]byte("Error to create statement"))
		return
	}
	defer statement.Close()

	//Executando o statement
	insert, error := statement.Exec(recipeType.Name, recipeType.Image)
	if error != nil {
		rw.Write([]byte("Error when running statement"))
		return
	}

	insertedId, error := insert.LastInsertId()
	if error != nil {
		rw.Write([]byte("Error when retrieve inserted id"))
		return
	}

	//status codes
	rw.WriteHeader(201)
	rw.Write([]byte(fmt.Sprintf("Type inserted with id: %d", insertedId)))

}

func ShowTypes(rw http.ResponseWriter, r *http.Request) {

	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close()

	lines, error := db.Query("select * from Types")
	if error != nil {
		rw.Write([]byte("Error when search types"))
		return
	}

	defer lines.Close()

	var recipeTypes []recipeType
	for lines.Next() {
		var recipeType recipeType

		if error := lines.Scan(&recipeType.ID, &recipeType.Name, &recipeType.Image); error != nil {
			rw.Write([]byte("Error when scanning recipeTypes"))
			return
		}

		recipeTypes = append(recipeTypes, recipeType)
	}

	rw.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(rw).Encode(recipeTypes); error != nil {
		rw.Write([]byte("Erro when convert recipeTypes to JSON"))
		return
	}

}
