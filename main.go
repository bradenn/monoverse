package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	start := time.Now()
	fmt.Printf("BIG BANG! [%s]\n", time.Since(start))
	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			return
		}
	}()
	monoverse := new(Monoverse)
	fmt.Printf("Created The Universe [%s]\n", time.Since(start))
	monoverse.Run()

}
