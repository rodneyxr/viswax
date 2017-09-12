package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// http://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item=21787
func main() {

	// Get the prices and profits for all runes
	prices := make(map[Rune]int)
	for _, r := range runeMap {
		r.Item = getItem(r.ID)
		prices[r] = r.GetProfit(40)
	}

	// Get rune combinations and each of their profits
	combos, err := getCombinations()
	if err != nil {
		panic(err)
	}

	type Pair struct {
		Key   Rune
		Value int
	}

	var pairs []Pair
	for k, v := range prices {
		pairs = append(pairs, Pair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	fmt.Println(combos)

	for _, kv := range pairs {
		fmt.Printf("    - %s:\t%4dK\n", strings.Replace(kv.Key.Name, " rune", "", 1), kv.Value/1000)
	}

	fmt.Println()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
