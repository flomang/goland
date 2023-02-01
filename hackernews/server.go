package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/flomang/goland/hackernews/graph"
	"github.com/flomang/goland/lib/postgres"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := postgres.InitDB("postgres://comms@localhost/hackernews?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	} else {
		defer db.Close()

		rows, _ := db.Query("SELECT username FROM users")
		var name string
		for rows.Next() {
			err := rows.Scan(&name)
			if err != nil {
				panic(err)
			}
			fmt.Println(name)
		}
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
