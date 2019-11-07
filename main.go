package main // import "github.com/etu/mkvcleaner"

import (
	"fmt"
	"github.com/giorgisio/goav/avformat"
)

func main() {
	fmt.Println("hello")

	// Register all formats and codecs
	avformat.AvRegisterAll()

	ctx := avformat.AvformatAllocContext()

	fmt.Println(ctx)
}
