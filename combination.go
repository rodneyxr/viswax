package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/yhat/scrape"
)

var viswax = getItem(32092)

// CombinationInfo contains information
type CombinationInfo struct {
	DayOfWeek  string
	DayOfMonth int
	Slot1      []SlotGroup
	Slot2      []SlotGroup
}

// SlotGroup contains a possible group of runes for slots
type SlotGroup struct {
	Group []Slot
}

// Slot contains information about runes and rewards
type Slot struct {
	RuneName string
	Count    int
	Rune     Rune
}

func jsonNumberToInt(number json.Number) int {
	sPrice := number.String()
	sPrice = strings.Replace(sPrice, ",", "", -1)
	price, _ := strconv.Atoi(sPrice)
	return price
}

func (c *CombinationInfo) String() string {
	sb := bytes.Buffer{}
	// Header
	sb.WriteString(fmt.Sprintf("Vis Wax for %s %d\n", c.DayOfWeek, c.DayOfMonth))
	sb.WriteString("-----------------------------\n")

	// Slot 1
	sb.WriteString("Slot 1:\n")
	for i, group := range c.Slot1 {
		sb.WriteString(fmt.Sprintf("  Group %d:\n", i+1))
		for _, slot := range group.Group {
			profit := (slot.Count*jsonNumberToInt(viswax.Current.Price) - slot.Rune.Amount*jsonNumberToInt(slot.Rune.Item.Current.Price)) / 1000
			sb.WriteString(fmt.Sprintf("    - %s:\t%d%10dK\n", strings.Title(slot.RuneName), slot.Count, profit))
		}
		sb.WriteRune('\n')
	}

	// Slot 2
	sb.WriteString("Slot 2:\n")
	for i, group := range c.Slot2 {
		sb.WriteString(fmt.Sprintf("  Group %d:\n", i+1))
		for _, slot := range group.Group {
			profit := (slot.Count*jsonNumberToInt(viswax.Current.Price) - slot.Rune.Amount*jsonNumberToInt(slot.Rune.Item.Current.Price)) / 1000
			sb.WriteString(fmt.Sprintf("    - %s:\t%d%10dK\n", strings.Title(slot.RuneName), slot.Count, profit))
		}
		sb.WriteRune('\n')
	}

	// Slot 3
	sb.WriteString("Slot 3:\n")
	sb.WriteString("    - Random:\t<=40\n")
	return sb.String()
}

func getCombinations() (*CombinationInfo, error) {
	combos := new(CombinationInfo)
	post, err := getPost()
	if err != nil {
		return nil, err
	}

	// Find the start of the combinations in the post
	// (Ex: Combination for Wednesday the 19th)
	re := regexp.MustCompile(`Combination\s*for\s*(\w+).*?(\d+)\w*;?`)
	res := re.FindStringSubmatch(post)

	combos.DayOfWeek = res[1]
	combos.DayOfMonth, err = strconv.Atoi(res[2])
	if err != nil {
		return nil, err
	}

	// Define regex for a combination group or slot label
	// Ex: "Air(earth 29)" or "Slot 1"
	re = regexp.MustCompile(`-\s*(\w+)|(\w+\*?)\s*(\d+)`)
	groups := re.FindAllStringSubmatch(post, -1)

	// Start with slot 0 and loop through each group. Each time a slot label
	// is found, increment the slot counter.
	slot := 0
	for _, group := range groups {
		var runeName string
		var count int

		// hasCount is true if there is more than one rune in the group
		hasCount := group[3] != ""

		if hasCount {
			runeName = strings.ToLower(group[2])
			if count, err = strconv.Atoi(group[3]); err != nil {
				return nil, err
			}
		} else {
			runeName = strings.ToLower(group[1])
		}

		if slot == 0 {
			// Slot 1
			if runeName == "slot" {
				slot = 1
				continue
			}
		} else if slot == 1 {
			// Slot 2
			if runeName == "slot" {
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
			key := strings.TrimRight(runeName, "*")
			key = strings.ToLower(key)
			slotToAdd := Slot{
				RuneName: runeName,
				Count:    count,
				Rune:     runeMap[key],
			}
			slotToAdd.Rune.Item = getItem(slotToAdd.Rune.ID)
			combos.Slot1[i].Group = append(combos.Slot1[i].Group, slotToAdd)

		} else if slot == 2 {
			// If the word "slot" or "is" is found, we found the "slot 3" label
			// or it is missing.
			if runeName == "slot" || runeName == "is" {
				break
			}
			if !hasCount {
				count = 30
			}
			if count == 30 {
				combos.Slot2 = append(combos.Slot2, SlotGroup{})
			}
			i := len(combos.Slot2) - 1
			key := strings.TrimRight(runeName, "*")
			key = strings.ToLower(key)
			slotToAdd := Slot{
				RuneName: runeName,
				Count:    count,
				Rune:     runeMap[key],
			}
			slotToAdd.Rune.Item = getItem(slotToAdd.Rune.ID)
			combos.Slot2[i].Group = append(combos.Slot2[i].Group, slotToAdd)
		}
	}
	return combos, nil
}

func getPost() (string, error) {
	// Get the forum page
	root, err := getHTML(`http://services.runescape.com/m=forum/forums.ws?75,76,331,66006366`)
	if err != nil {
		return "", err
	}

	// Get all posts from forum
	postMatcher := scrape.ByClass("forum-post__body")
	nodes := scrape.FindAll(root, postMatcher)
	if len(nodes) < 2 {
		fmt.Fprintln(os.Stderr, "Failed to find posts.")
		return "", err
	}

	// First post is the combinations post
	return scrape.Text(nodes[0]), nil
}
