package main

import (
	"taskmanagement/Delivery/router"
)

func main() {

	router := router.Router()
	router.Run(":3030")
}
