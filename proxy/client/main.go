package main

import (
	"fmt"
)

func main(){
	exitCh := make(chan error, 1)

	setPAC()
	fmt.Println("SetPaCOK")

	<-exitCh
}
