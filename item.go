package main

import (
	"encoding/json"
	"fmt"
)

var itemCache = make(map[int]Item)

type itemWrapper struct {
	Item Item `json:"item"`
}

// Item conatains information about an item returned from the itemdb api
type Item struct {
	Icon        string `json:"icon"`
	IconLarge   string `json:"icon_large"`
	ID          int    `json:"id"`
	Type        string `json:"type"`
	TypeIcon    string `json:"typeIcon"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsMembers   string `json:"members"`
	Current     Price  `json:"current"`
}

// Price holds the price and trend of the item
type Price struct {
	Trend string      `json:"trend"`
	Price json.Number `json:"price,Number"`
}

func getItem(id int) Item {
	if item, ok := itemCache[id]; ok {
		return item
	}
	url := fmt.Sprintf("http://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item=%d", id)
	wrapper := &itemWrapper{}
	getJSON(url, wrapper)
	itemCache[id] = wrapper.Item
	return wrapper.Item
}
