package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	// "testTech/src/routes"
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
	fmt.Println("Server run on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
	fmt.Println("Exit server")
	routes.GetDomainURL()
}