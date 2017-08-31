package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

// CombinationInfo contains information
type CombinationInfo struct {
	Slot1 []SlotGroup
	Slot2 []SlotGroup
}

// SlotGroup contains a possible group of runes for slots
type SlotGroup struct {
	Group []Slot
}

// Slot contains information about runes and rewards
type Slot struct {
	Rune  string
	Count int
}

func main() {
	// Get the forum post
	root, err := getHTML(`http://services.runescape.com/m=forum/forums.ws?75,76,387,65763383`)
	if err != nil {
		panic(err)
	}

	// Get all posts from forum
	postMatcher := scrape.ByClass("forum-post__body")
	nodes := scrape.FindAll(root, postMatcher)
	if len(nodes) < 2 {
		fmt.Fprintln(os.Stderr, "Failed to find posts.")
		os.Exit(1)
	}

	// Second post is the combinations post
	post := scrape.Text(nodes[1])

	re, err := regexp.Compile(`Combination\s*for\s*(\w+).*?(\d+)\w*;?`)
	if err != nil {
		panic(err)
	}
	res := re.FindStringSubmatch(post)
	fmt.Printf("Vis Wax for %s %s\n", res[1], res[2])
	fmt.Println("-----------------------------")

	combos, err := getCombinations(post)
	if err != nil {
		panic(err)
	}

	// Slot 1
	fmt.Println("Slot 1:")
	for i, group := range combos.Slot1 {
		fmt.Printf("  Group %d:\n", i+1)
		for _, slot := range group.Group {
			fmt.Printf("    - %s:\t%d\n", strings.Title(slot.Rune), slot.Count)
		}
		fmt.Println()
	}

	fmt.Println()

	// Slot 2
	fmt.Println("Slot 2:")
	for i, group := range combos.Slot2 {
		fmt.Printf("  Group %d:\n", i+1)
		for _, slot := range group.Group {
			fmt.Printf("    - %s:\t%d\n", strings.Title(slot.Rune), slot.Count)
		}
		fmt.Println()
	}
}

func getHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return root, nil
}

func getCombinations(post string) (*CombinationInfo, error) {
	re, err := regexp.Compile(`-\s*(\w+)|(\w+\*?)\s*(\d+)`)
	if err != nil {
		return nil, err
	}
	groups := re.FindAllStringSubmatch(post, -1)
	combos := new(CombinationInfo)
	slot := 0
	for _, group := range groups {
		var rune string
		var count int
		hasCount := group[3] != ""

		if hasCount {
			rune = strings.ToLower(group[2])
			if count, err = strconv.Atoi(group[3]); err != nil {
				return nil, err
			}
		} else {
			rune = strings.ToLower(group[1])
		}

		if slot == 0 {
			if rune == "slot" {
				slot = 1
				continue
			}
		} else if slot == 1 {
			if rune == "slot" {
				slot = 2
				continue
			}
			if !hasCount {
				count = 30
			}
			if count == 30 {
				combos.Slot1 = append(combos.Slot1, SlotGroup{})
			}
			i := len(combos.Slot1) - 1
			combos.Slot1[i].Group = append(combos.Slot1[i].Group, Slot{
				Rune:  rune,
				Count: count,
			})
		} else if slot == 2 {
			if rune == "slot" {
				break
			}
			if !hasCount {
				count = 30
			}
			if count == 30 {
				combos.Slot2 = append(combos.Slot2, SlotGroup{})
			}
			i := len(combos.Slot2) - 1
			combos.Slot2[i].Group = append(combos.Slot2[i].Group, Slot{
				Rune:  rune,
				Count: count,
			})
		}
	}
	return combos, nil
}
