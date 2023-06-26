package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("%d", time.Now().UnixNano()/1000000000)
}
