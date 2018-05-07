package main

import (
	"flag"
	"time"

)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	after := time.After(3*time.Second)
	<-after
	println("==============")
}