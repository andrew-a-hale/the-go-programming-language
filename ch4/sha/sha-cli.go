// take stdin and return sha256 by default and sha384 and sha512 by flag args
package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"log"
	"os"
)

var shaType *string = flag.String("sha", "256", "types: {256, 384, 512} default is 256")

func main() {
	flag.Parse()

	buf := bufio.NewReader(os.Stdin)
	input, err := buf.ReadBytes('\n')
	if err != nil {
		log.Fatal("could not read stdin")
	}

	switch *shaType {
	case "512":
		fmt.Printf("%x\n", sha512.Sum512(input))
	case "384":
		fmt.Printf("%x\n", sha512.Sum384(input))
	case "256":
		fmt.Printf("%x\n", sha256.Sum256(input))
	default:
		log.Fatalln("invalid sha type")
	}
}
