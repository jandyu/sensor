package main

import "fmt"

//var addr = flag.String("addr", ":8901", "http service address")

func main() {

	opt := &ServerOpts{
		Name: "Monitor",
		Port: 9000,
	}
	// start server
	mserver := NewServer(opt)
	err := mserver.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
