package main

import "fmt"

var version = "LOCAL_BUILD" // Gets changed by CI system

func main() {
	fmt.Println(version)
}
