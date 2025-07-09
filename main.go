package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")

}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	splitText := strings.Fields(lowerText)
	return splitText
}
