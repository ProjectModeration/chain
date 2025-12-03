package main

import (
	inter "ProjectModeration/chain/interaction"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	res := inter.TextClassification("hello!")

	fmt.Println("Response from interaction:", res)

	res2 := inter.ImageClassification("noFilter.webp")
	fmt.Println("Image Classification Response:", res2)
}
