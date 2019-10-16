package main

import (
	"log"
	"net/http"

	yss "yss-go-official/api"
)

func main() {
	err := http.ListenAndServe(":8080", yss.YssRouter)
	if err != nil {
		log.Fatal("Error at listen", err)
	}
}
