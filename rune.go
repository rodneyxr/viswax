package main

// Rune describes the id, name and price of a rune.
type Rune struct {
	ID     int
	Name   string
	Amount int
	Item   Item
}

var runeMap = map[string]Rune{
	// 300
	"astral": Rune{ID: 9075, Name: "Astral rune", Amount: 300},
	"soul":   Rune{ID: 566, Name: "Soul rune", Amount: 300},
	"law":    Rune{ID: 563, Name: "Law rune", Amount: 300},
	"mud":    Rune{ID: 4698, Name: "Mud rune", Amount: 300},

	// 350
	"blood":  Rune{ID: 565, Name: "Blood rune", Amount: 350},
	"nature": Rune{ID: 561, Name: "Nature rune", Amount: 350},

	// 400
	"cosmic": Rune{ID: 564, Name: "Cosmic rune", Amount: 400},
	"death":  Rune{ID: 560, Name: "Death rune", Amount: 400},

	// 500
	"lava":  Rune{ID: 4699, Name: "Lava rune", Amount: 500},
	"steam": Rune{ID: 4694, Name: "Steam rune", Amount: 500},
	"chaos": Rune{ID: 562, Name: "Chaos rune", Amount: 500},
	"smoke": Rune{ID: 4697, Name: "Smoke rune", Amount: 500},
	"dust":  Rune{ID: 4696, Name: "Dust rune", Amount: 500},
	"mist":  Rune{ID: 4695, Name: "Mist rune", Amount: 500},

	// 1000
	"fire":  Rune{ID: 554, Name: "Fire rune", Amount: 1000},
	"earth": Rune{ID: 557, Name: "Earth rune", Amount: 1000},
	"water": Rune{ID: 555, Name: "Water rune", Amount: 1000},
	"air":   Rune{ID: 556, Name: "Air rune", Amount: 1000},

	// 2000
	"mind": Rune{ID: 558, Name: "Mind rune", Amount: 2000},
	"body": Rune{ID: 559, Name: "Body rune", Amount: 2000},
}

// GetProfit returns the profit of using the rune given the number of viswax returned
func (r *Rune) GetProfit(count int) int {
	viswaxPrice := count * jsonNumberToInt(viswax.Current.Price)
	runePrice := r.Amount * jsonNumberToInt(r.Item.Current.Price)
	return viswaxPrice - runePrice
}
