package main

import (
	"fmt"
	_ "net/http/pprof"
	"time"
)

func main() {
	start := time.Now()
	fmt.Printf("BIG BANG! [%s]\n", time.Since(start))
	// go http.ListenAndServe("localhost:8080", nil)
	monoverse := new(Monoverse)
	fmt.Printf("Created The Universe [%s]\n", time.Since(start))
	monoverse.Run()

}
