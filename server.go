package main

import (
	"log"
	"net/http"
)

func networkError(err error, c chan struct{}) {
	log.Println(err)
	if err.Error() == "use of closed network connection" {
		close(c)
	}
}

func modelError(err error, c chan struct{}) {
	log.Println(err)
	// do something with the error value
	close(c)
}

func dataError(err error, c chan struct{}) {
	log.Println(err)
	// do something with the error value
	close(c)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
