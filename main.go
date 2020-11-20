package main

import (
	"fmt"
	"sync"
)

//var addr = flag.String("addr", ":8901", "http service address")

var MDATA sync.Map

func main() {

	//初始化
	for i := 1; i < 10; i++ {
		MDATA.Store("LED"+fmt.Sprint(i)+"_tmp_m", "0.00")
		MDATA.Store("LED"+fmt.Sprint(i)+"_sal_m", "0.00")
	}

	go startHttp(8081)

	// 112.16.93.184 9123
	opt := &ServerOpts{
		Name: "Monitor",
		Port: 8901,
	}
	// start server
	mserver := NewServer(opt)
	err := mserver.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
