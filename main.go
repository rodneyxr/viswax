package main

import (
	"bufio"
	"fmt"
	"os"
)

// http://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item=21787

func main() {

	combos, err := getCombinations()
	if err != nil {
		panic(err)
	}

	fmt.Println(combos)
	fmt.Println()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
