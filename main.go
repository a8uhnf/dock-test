package main

import (
	"github.com/a8uhnf/container-h/pkg/controller"
)

func main() {
	go controller.CreateInformerQueueNs()
	// go controller.CreateInformerQueue()
	select {}
}
