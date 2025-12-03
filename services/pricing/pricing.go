package main

// Grille simple A → B
var pricingTable = map[string]float64{
	"Centre-ville:Aéroport": 25.5,
	"Plateau:McGill":        12.0,
	"Aéroport:Centre-ville": 25.5,
	"McGill:Plateau":        12.0,
}

// Récupère le prix des zones
func GetPrice(from string, to string) (float64, bool) {
	key := from + ":" + to
	value, exists := pricingTable[key]
	return value, exists
}
