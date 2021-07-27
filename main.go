package main

import (
	"fmt"

	pkgApplication "github.com/tmtriet200800/test-go-utils/application"
)

func main() {
	fmt.Println("Welcome to our internal package")

	pkgApplication.Help()
}