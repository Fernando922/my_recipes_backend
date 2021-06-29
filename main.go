package main

import (
	"fmt"
	"log"
	"net/http"
	"recipes-backend/server"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//Categorias
	router.HandleFunc("/categories", server.CreateCategory).Methods(http.MethodPost)
	router.HandleFunc("/categories", server.ShowCategories).Methods(http.MethodGet)
	router.HandleFunc("/categories/{id}", server.UpdateCategory).Methods(http.MethodPut)
	router.HandleFunc("/categories/{id}", server.DeleteCategory).Methods(http.MethodDelete)

	//Receitas
	router.HandleFunc("/recipes", server.CreateRecipe).Methods(http.MethodPost)
	router.HandleFunc("/recipes/category/{id}", server.ShowRecipes).Methods(http.MethodGet)
	router.HandleFunc("/recipes/{id}", server.UpdateRecipe).Methods(http.MethodPut)
	router.HandleFunc("/recipes/{id}", server.DeleteRecipe).Methods(http.MethodDelete)

	//Ingredientes
	router.HandleFunc("/ingredients", server.CreateIngredient).Methods(http.MethodPost)
	router.HandleFunc("/ingredients/recipe/{id}", server.ShowIngredients).Methods(http.MethodGet)
	router.HandleFunc("/ingredients/{id}", server.UpdateIngredient).Methods(http.MethodPut)
	router.HandleFunc("/ingredients/{id}", server.DeleteIngredient).Methods(http.MethodDelete)

	//Dicas
	router.HandleFunc("/tips", server.CreateTip).Methods(http.MethodPost)
	router.HandleFunc("/tips/recipe/{id}", server.ShowTips).Methods(http.MethodGet)
	router.HandleFunc("/tips/{id}", server.UpdateTip).Methods(http.MethodPut)
	router.HandleFunc("/tips/{id}", server.DeleteTip).Methods(http.MethodDelete)

	fmt.Println("Server listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
