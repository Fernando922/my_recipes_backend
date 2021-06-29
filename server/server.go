package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"recipes-backend/database"
	"strconv"

	"github.com/gorilla/mux"
)

type recipe struct {
	ID           uint32 `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	CategoryId   uint32 `json:"categoryId`
	Prepare_mode string `json:"prepare_mode"`
}

type category struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type ingredient struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	RecipeId uint32 `json:"recipeId"`
}

type tip struct {
	ID          uint32 `json:"id"`
	Description string `json:"description"`
	RecipeId    uint32 `json:"recipeId`
}

func CreateCategory(rw http.ResponseWriter, r *http.Request) {

	//lendo o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Request body read error"))
		return
	}

	//convertendo o json do body da requisição para struct
	var category category
	if error = json.Unmarshal(requestBody, &category); error != nil {
		rw.Write([]byte("Error to convert type to struct"))
		return
	}

	//abrindo conexão com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close() //só será executado no final da função

	//Criando o prepare statement (comando de inserção, para evitar sql injection)
	statement, error := db.Prepare("insert into Categories (name, image) values (?,?)")
	if error != nil {
		rw.Write([]byte("Error to create statement"))
		return
	}
	defer statement.Close()

	//Executando o statement
	insert, error := statement.Exec(category.Name, category.Image)
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

func ShowCategories(rw http.ResponseWriter, r *http.Request) {

	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close()

	lines, error := db.Query("select * from Categories")
	if error != nil {
		rw.Write([]byte("Error when search types"))
		return
	}

	defer lines.Close()

	var categories []category
	for lines.Next() {
		var category category

		if error := lines.Scan(&category.ID, &category.Name, &category.Image); error != nil {
			rw.Write([]byte("Error when scanning categorys"))
			return
		}

		categories = append(categories, category)
	}

	rw.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(rw).Encode(categories); error != nil {
		rw.Write([]byte("Erro when convert categorys to JSON"))
		return
	}

}

func UpdateCategory(rw http.ResponseWriter, r *http.Request) {

	//Pegando os parametros da requisição
	parameters := mux.Vars(r) //pega os params, é retornado um MAP de string
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)

	if error != nil {
		rw.Write([]byte("Error while convert parameter to integer"))
		return
	}

	//Pegando o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Error when reading request body"))
		return
	}

	//Convertendo para struct
	var category category
	if error := json.Unmarshal(requestBody, &category); error != nil {
		rw.Write([]byte("Erro converting category to struct"))
	}

	//Conectando com o banco de dados
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Database connect error"))
		return
	}
	defer db.Close()

	//criando um statement
	statement, error := db.Prepare("update Categories set name=?, image=? where id = ?")
	if error != nil {
		rw.Write([]byte("Error creating statement"))
		return
	}
	defer statement.Close()

	if _, error := statement.Exec(category.Name, category.Image, ID); error != nil {
		rw.Write([]byte("Error updating category"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func DeleteCategory(rw http.ResponseWriter, r *http.Request) {
	//Pegando os parametros da requisição
	parameters := mux.Vars(r) //pega os params, é retornado um MAP de string
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)

	if error != nil {
		rw.Write([]byte("Error while convert parameter to integer"))
		return
	}

	//Conectando com o banco de dados
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Database connect error"))
		return
	}
	defer db.Close()

	//criando um statement
	statement, error := db.Prepare("delete from Categories where id = ?")
	if error != nil {
		rw.Write([]byte("Error creating statement"))
		return
	}
	defer statement.Close()

	if _, error := statement.Exec(ID); error != nil {
		rw.Write([]byte("Error updating category"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func CreateRecipe(rw http.ResponseWriter, r *http.Request) {
	//lendo o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Request body read error"))
		return
	}

	//convertendo o json do body da requisição para struct
	var recipe recipe
	if error = json.Unmarshal(requestBody, &recipe); error != nil {
		rw.Write([]byte("Error to convert type to struct"))
		return
	}

	//abrindo conexão com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close() //só será executado no final da função

	//Criando o prepare statement (comando de inserção, para evitar sql injection)
	statement, error := db.Prepare("insert into Recipes (name, description, image, categoryId, prepare_mode) values (?,?,?,?,?)")
	if error != nil {
		rw.Write([]byte("Error to create statement"))
		return
	}
	defer statement.Close()

	//Executando o statement
	insert, error := statement.Exec(recipe.Name, recipe.Description, recipe.Image, recipe.CategoryId, recipe.Prepare_mode)
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
	rw.Write([]byte(fmt.Sprintf("Recipe inserted with id: %d", insertedId)))

}

func ShowRecipes(rw http.ResponseWriter, r *http.Request) {

	//Pegando os parametros da requisição
	parameters := mux.Vars(r) //pega os params, é retornado um MAP de string
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)

	if error != nil {
		rw.Write([]byte("Error while convert parameter to integer"))
		return
	}

	//Conectar com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close()

	//Buscar todos os registros no banco
	lines, error := db.Query("select * from Recipes where categoryId=?", ID)
	if error != nil {
		rw.Write([]byte("Error when search recipes"))
		return
	}

	defer lines.Close()

	var recipes []recipe
	for lines.Next() {
		var recipe recipe

		if error := lines.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.Image, &recipe.CategoryId, &recipe.Prepare_mode); error != nil {
			rw.Write([]byte("Error when scanning recipes"))
			return
		}

		recipes = append(recipes, recipe)
	}

	rw.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(rw).Encode(recipes); error != nil {
		rw.Write([]byte("Erro when convert categorys to JSON"))
		return
	}
}

func UpdateRecipe(rw http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)
	if error != nil {
		rw.Write([]byte("Error when convert parameter to integer"))
		return
	}

	//Pegando o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Error when read request body"))
		return
	}

	//Convertendo o corpo da requisição para um struct
	var recipe recipe
	if error := json.Unmarshal(requestBody, &recipe); error != nil {
		rw.Write([]byte("Error when convert recipe to struct"))
		return
	}

	//conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}

	defer db.Close()

	//criando o statement para fazer a alteração no banco
	statement, error := db.Prepare("update Recipes set name = ?, description = ?, image = ?, categoryId = ?, prepare_mode = ? where id = ?")
	if error != nil {
		rw.Write([]byte("Error when create statement!"))
		return
	}
	defer statement.Close()

	//executando o comando sql do statement criado anteriormente
	if _, error := statement.Exec(recipe.Name, recipe.Description, recipe.Image, recipe.CategoryId, recipe.Prepare_mode, ID); error != nil {
		rw.Write([]byte("Error when update recipe"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func DeleteRecipe(rw http.ResponseWriter, r *http.Request) {
	//recuperando o parametro passado na requisição
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)
	if error != nil {
		rw.Write([]byte("Error to convter parameter to integer"))
		return
	}

	//conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}

	defer db.Close()

	//criando o statement para executar no banco
	statement, error := db.Prepare("delete from Recipes where id = ?")
	if error != nil {
		rw.Write([]byte("Error when create statement"))
	}
	defer statement.Close()

	//executando o statement
	if _, error := statement.Exec(ID); error != nil {
		rw.Write([]byte("Error to delete recipe"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func CreateIngredient(rw http.ResponseWriter, r *http.Request) {

	//lendo o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Request body read eror"))
	}

	//convertendo o json do boyd da requisição para struct
	var ingredient ingredient
	if error := json.Unmarshal(requestBody, &ingredient); error != nil {
		rw.Write([]byte("Error to convert ingredient to struct"))
		return
	}

	//abrindo conexão com obanco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close()

	//criando o statement para manipular o banco de dados
	statement, error := db.Prepare("insert into Ingredients (name, recipeId) values (?,?)")
	if error != nil {
		rw.Write([]byte("Error to create statement"))
		return
	}

	defer statement.Close()

	//Executando o statement
	insert, error := statement.Exec(ingredient.Name, ingredient.RecipeId)
	if error != nil {
		rw.Write([]byte("Error when running statement"))
		return
	}

	//recuperando o id do novo registro no banco
	insertedId, error := insert.LastInsertId()
	if error != nil {
		rw.Write([]byte("Erro when retrieve inserted id"))
		return
	}

	//status codes
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fmt.Sprintf("Ingredient inserted with id: %d", insertedId)))

}

func ShowIngredients(rw http.ResponseWriter, r *http.Request) {

	//Pegando os paramentros da requisição
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)

	if error != nil {
		rw.Write([]byte("Error while convert parameter to integer"))
		return
	}

	//Conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}

	defer db.Close()

	//Buscar todos os registros no banco
	lines, error := db.Query("select * from Ingredients where recipeId=?", ID)
	if error != nil {
		rw.Write([]byte("Error when search ingredients"))
		return
	}
	defer lines.Close()

	var ingredients []ingredient
	for lines.Next() {
		var ingredient ingredient
		if error := lines.Scan(&ingredient.ID, &ingredient.Name, &ingredient.RecipeId); error != nil {
			rw.Write([]byte("Error when search ingredients"))
			return
		}

		ingredients = append(ingredients, ingredient)
	}

	rw.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(rw).Encode(ingredients); error != nil {
		rw.Write([]byte("Error when convert ingredients to JSON"))
		return
	}

}

func UpdateIngredient(rw http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)
	if error != nil {
		rw.Write([]byte("Error when convert parameter to integer"))
		return
	}

	//pegando o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Error when read request body"))
		return
	}

	var ingredient ingredient
	if error := json.Unmarshal(requestBody, &ingredient); error != nil {
		rw.Write([]byte("Error when convert ingredients to struct"))
		return
	}

	//conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close()

	//criando o statement para fazer a alteração no banco
	statement, error := db.Prepare("update Ingredients set name=?, recipeId = ? where id = ?")
	if error != nil {
		rw.Write([]byte("Error when create statement"))
		return
	}
	defer statement.Close()

	//executando o comando sql criado anteriormente
	if _, error := statement.Exec(ingredient.Name, ingredient.RecipeId, ID); error != nil {
		fmt.Println(error)
		rw.Write([]byte("Error when update ingredient"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func DeleteIngredient(rw http.ResponseWriter, r *http.Request) {

	//recuperando o parametro passado na requisição
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)
	if error != nil {
		rw.Write([]byte("Error to convert parameter to integer"))
		return
	}

	//conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to convert parameter to integer"))
		return
	}
	defer db.Close()

	//criando o statement para executar no banco
	statement, error := db.Prepare("delete from Ingredients where id = ?")
	if error != nil {
		rw.Write([]byte("Error when create statement"))
		return
	}
	defer statement.Close()

	//executando o statement
	if _, error := statement.Exec(ID); error != nil {
		rw.Write([]byte("Error to delete ingredient"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func CreateTip(rw http.ResponseWriter, r *http.Request) {

	//lendo o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Request body read eror"))
	}

	//convertendo o json do boyd da requisição para struct
	var tip tip
	if error := json.Unmarshal(requestBody, &tip); error != nil {
		fmt.Println(error)
		rw.Write([]byte("Error to convert tip to struct"))
		return
	}

	//abrindo conexão com obanco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close()

	//criando o statement para manipular o banco de dados
	statement, error := db.Prepare("insert into Tips (description, recipeId) values (?,?)")
	if error != nil {
		rw.Write([]byte("Error to create statement"))
		return
	}

	defer statement.Close()

	//Executando o statement
	insert, error := statement.Exec(tip.Description, tip.RecipeId)
	if error != nil {
		rw.Write([]byte("Error when running statement"))
		return
	}

	//recuperando o id do novo registro no banco
	insertedId, error := insert.LastInsertId()
	if error != nil {
		rw.Write([]byte("Erro when retrieve inserted id"))
		return
	}

	//status codes
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fmt.Sprintf("Tip inserted with id: %d", insertedId)))

}

func ShowTips(rw http.ResponseWriter, r *http.Request) {

	//Pegando os paramentros da requisição
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)

	if error != nil {
		rw.Write([]byte("Error while convert parameter to integer"))
		return
	}

	//Conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}

	defer db.Close()

	//Buscar todos os registros no banco
	lines, error := db.Query("select * from Tips where recipeId=?", ID)
	if error != nil {
		rw.Write([]byte("Error when search ingredients"))
		return
	}
	defer lines.Close()

	var tips []tip
	for lines.Next() {
		var tip tip
		if error := lines.Scan(&tip.ID, &tip.Description, &tip.RecipeId); error != nil {
			rw.Write([]byte("Error when search ingredients"))
			return
		}

		tips = append(tips, tip)
	}

	rw.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(rw).Encode(tips); error != nil {
		rw.Write([]byte("Error when convert ingredients to JSON"))
		return
	}

}

func UpdateTip(rw http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)
	if error != nil {
		rw.Write([]byte("Error when convert parameter to integer"))
		return
	}

	//pegando o corpo da requisição
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		rw.Write([]byte("Error when read request body"))
		return
	}

	var tip tip
	if error := json.Unmarshal(requestBody, &tip); error != nil {
		rw.Write([]byte("Error when convert tips to struct"))
		return
	}

	//conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to connect to database"))
		return
	}
	defer db.Close()

	//criando o statement para fazer a alteração no banco
	statement, error := db.Prepare("update Tips set description=?, recipeId = ? where id = ?")
	if error != nil {
		rw.Write([]byte("Error when create statement"))
		return
	}
	defer statement.Close()

	//executando o comando sql criado anteriormente
	if _, error := statement.Exec(tip.Description, tip.RecipeId, ID); error != nil {
		fmt.Println(error)
		rw.Write([]byte("Error when update tip"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func DeleteTip(rw http.ResponseWriter, r *http.Request) {

	//recuperando o parametro passado na requisição
	parameters := mux.Vars(r)
	ID, error := strconv.ParseInt(parameters["id"], 10, 32)
	if error != nil {
		rw.Write([]byte("Error to convert parameter to integer"))
		return
	}

	//conectando com o banco
	db, error := database.Connect()
	if error != nil {
		rw.Write([]byte("Error to convert parameter to integer"))
		return
	}
	defer db.Close()

	//criando o statement para executar no banco
	statement, error := db.Prepare("delete from Tips where id = ?")
	if error != nil {
		rw.Write([]byte("Error when create statement"))
		return
	}
	defer statement.Close()

	//executando o statement
	if _, error := statement.Exec(ID); error != nil {
		rw.Write([]byte("Error to delete tip"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}
