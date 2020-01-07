package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"apiDomainInfo/routes"
)

//Routes return routes
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer, middleware.Logger)
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/servers", routes.Server())
	})
	return router
}

func main() {
	router := Routes()
	vars := os.Environ()
	fmt.Println("las variables", vars)
	fmt.Println("Server run on port",os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+ os.Getenv("PORT") , router))
	fmt.Println("Exit server")
	// routes.GetDomainURL()
}