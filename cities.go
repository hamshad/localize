package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// City represents a city with timezone and location information.
type City struct {
	Name        string
	Timezone    string
	Country     string
	Category    string     // Americas, Europe, MiddleEast, Asia, Africa, Oceania
	Coordinates [2]float64 // Latitude, Longitude (for future map markers)
	Color       tcell.Color
}

// AllCities is the comprehensive database of cities.
var AllCities = []City{
	// Americas
	{Name: "Honolulu", Timezone: "Pacific/Honolulu", Country: "USA", Category: "Americas", Coordinates: [2]float64{21.3069, -157.8583}, Color: tcell.ColorTurquoise},
	{Name: "Anchorage", Timezone: "America/Anchorage", Country: "USA", Category: "Americas", Coordinates: [2]float64{61.2181, -149.9003}, Color: tcell.ColorLightCyan},
	{Name: "Los Angeles", Timezone: "America/Los_Angeles", Country: "USA", Category: "Americas", Coordinates: [2]float64{34.0522, -118.2437}, Color: tcell.ColorDarkCyan},
	{Name: "Vancouver", Timezone: "America/Vancouver", Country: "Canada", Category: "Americas", Coordinates: [2]float64{49.2827, -123.1207}, Color: tcell.ColorCadetBlue},
	{Name: "San Francisco", Timezone: "America/Los_Angeles", Country: "USA", Category: "Americas", Coordinates: [2]float64{37.7749, -122.4194}, Color: tcell.ColorSteelBlue},
	{Name: "Denver", Timezone: "America/Denver", Country: "USA", Category: "Americas", Coordinates: [2]float64{39.7392, -104.9903}, Color: tcell.ColorMediumSeaGreen},
	{Name: "Chicago", Timezone: "America/Chicago", Country: "USA", Category: "Americas", Coordinates: [2]float64{41.8781, -87.6298}, Color: tcell.ColorSeaGreen},
	{Name: "Toronto", Timezone: "America/Toronto", Country: "Canada", Category: "Americas", Coordinates: [2]float64{43.6532, -79.3832}, Color: tcell.ColorDarkSlateGray},
	{Name: "New York", Timezone: "America/New_York", Country: "USA", Category: "Americas", Coordinates: [2]float64{40.7128, -74.0060}, Color: tcell.ColorDodgerBlue},
	{Name: "Miami", Timezone: "America/New_York", Country: "USA", Category: "Americas", Coordinates: [2]float64{25.7617, -80.1918}, Color: tcell.ColorDeepSkyBlue},
	{Name: "Mexico City", Timezone: "America/Mexico_City", Country: "Mexico", Category: "Americas", Coordinates: [2]float64{19.4326, -99.1332}, Color: tcell.ColorLawnGreen},
	{Name: "Sao Paulo", Timezone: "America/Sao_Paulo", Country: "Brazil", Category: "Americas", Coordinates: [2]float64{-23.5505, -46.6333}, Color: tcell.ColorLimeGreen},
	{Name: "Buenos Aires", Timezone: "America/Argentina/Buenos_Aires", Country: "Argentina", Category: "Americas", Coordinates: [2]float64{-34.6037, -58.3816}, Color: tcell.ColorMediumSpringGreen},

	// Europe
	{Name: "London", Timezone: "Europe/London", Country: "UK", Category: "Europe", Coordinates: [2]float64{51.5074, -0.1278}, Color: tcell.ColorGreen},
	{Name: "Lisbon", Timezone: "Europe/Lisbon", Country: "Portugal", Category: "Europe", Coordinates: [2]float64{38.7223, -9.1393}, Color: tcell.ColorPaleGreen},
	{Name: "Amsterdam", Timezone: "Europe/Amsterdam", Country: "Netherlands", Category: "Europe", Coordinates: [2]float64{52.3676, 4.9041}, Color: tcell.ColorSpringGreen},
	{Name: "Paris", Timezone: "Europe/Paris", Country: "France", Category: "Europe", Coordinates: [2]float64{48.8566, 2.3522}, Color: tcell.ColorDarkGreen},
	{Name: "Berlin", Timezone: "Europe/Berlin", Country: "Germany", Category: "Europe", Coordinates: [2]float64{52.5200, 13.4050}, Color: tcell.ColorForestGreen},
	{Name: "Stockholm", Timezone: "Europe/Stockholm", Country: "Sweden", Category: "Europe", Coordinates: [2]float64{59.3293, 18.0686}, Color: tcell.ColorYellowGreen},
	{Name: "Warsaw", Timezone: "Europe/Warsaw", Country: "Poland", Category: "Europe", Coordinates: [2]float64{52.2297, 21.0122}, Color: tcell.ColorDarkSeaGreen},
	{Name: "Athens", Timezone: "Europe/Athens", Country: "Greece", Category: "Europe", Coordinates: [2]float64{37.9838, 23.7275}, Color: tcell.ColorMediumAquamarine},
	{Name: "Moscow", Timezone: "Europe/Moscow", Country: "Russia", Category: "Europe", Coordinates: [2]float64{55.7558, 37.6173}, Color: tcell.ColorRed},
	{Name: "Istanbul", Timezone: "Europe/Istanbul", Country: "Turkey", Category: "Europe", Coordinates: [2]float64{41.0082, 28.9784}, Color: tcell.ColorIndianRed},

	// Middle East
	{Name: "Dubai", Timezone: "Asia/Dubai", Country: "UAE", Category: "MiddleEast", Coordinates: [2]float64{25.2048, 55.2708}, Color: tcell.ColorGold},
	{Name: "Tel Aviv", Timezone: "Asia/Jerusalem", Country: "Israel", Category: "MiddleEast", Coordinates: [2]float64{32.0853, 34.7818}, Color: tcell.ColorOrangeRed},
	{Name: "Riyadh", Timezone: "Asia/Riyadh", Country: "Saudi Arabia", Category: "MiddleEast", Coordinates: [2]float64{24.7136, 46.6753}, Color: tcell.ColorDarkOrange},

	// Africa
	{Name: "Cairo", Timezone: "Africa/Cairo", Country: "Egypt", Category: "Africa", Coordinates: [2]float64{30.0444, 31.2357}, Color: tcell.ColorSandyBrown},
	{Name: "Lagos", Timezone: "Africa/Lagos", Country: "Nigeria", Category: "Africa", Coordinates: [2]float64{6.5244, 3.3792}, Color: tcell.ColorChocolate},
	{Name: "Johannesburg", Timezone: "Africa/Johannesburg", Country: "South Africa", Category: "Africa", Coordinates: [2]float64{-26.2041, 28.0473}, Color: tcell.ColorPeru},
	{Name: "Nairobi", Timezone: "Africa/Nairobi", Country: "Kenya", Category: "Africa", Coordinates: [2]float64{-1.2921, 36.8219}, Color: tcell.ColorCoral},

	// Asia
	{Name: "Karachi", Timezone: "Asia/Karachi", Country: "Pakistan", Category: "Asia", Coordinates: [2]float64{24.8607, 67.0011}, Color: tcell.ColorDarkSalmon},
	{Name: "Mumbai", Timezone: "Asia/Kolkata", Country: "India", Category: "Asia", Coordinates: [2]float64{19.0760, 72.8777}, Color: tcell.ColorOrange},
	{Name: "Bangkok", Timezone: "Asia/Bangkok", Country: "Thailand", Category: "Asia", Coordinates: [2]float64{13.7563, 100.5018}, Color: tcell.ColorTomato},
	{Name: "Jakarta", Timezone: "Asia/Jakarta", Country: "Indonesia", Category: "Asia", Coordinates: [2]float64{-6.2088, 106.8456}, Color: tcell.ColorCrimson},
	{Name: "Singapore", Timezone: "Asia/Singapore", Country: "Singapore", Category: "Asia", Coordinates: [2]float64{1.3521, 103.8198}, Color: tcell.ColorDarkMagenta},
	{Name: "Manila", Timezone: "Asia/Manila", Country: "Philippines", Category: "Asia", Coordinates: [2]float64{14.5995, 120.9842}, Color: tcell.ColorHotPink},
	{Name: "Hong Kong", Timezone: "Asia/Hong_Kong", Country: "Hong Kong", Category: "Asia", Coordinates: [2]float64{22.3193, 114.1694}, Color: tcell.ColorDarkMagenta},
	{Name: "Shanghai", Timezone: "Asia/Shanghai", Country: "China", Category: "Asia", Coordinates: [2]float64{31.2304, 121.4737}, Color: tcell.ColorOrangeRed},
	{Name: "Seoul", Timezone: "Asia/Seoul", Country: "South Korea", Category: "Asia", Coordinates: [2]float64{37.5665, 126.9780}, Color: tcell.ColorMediumVioletRed},
	{Name: "Tokyo", Timezone: "Asia/Tokyo", Country: "Japan", Category: "Asia", Coordinates: [2]float64{35.6762, 139.6503}, Color: tcell.ColorDeepPink},

	// Oceania
	{Name: "Auckland", Timezone: "Pacific/Auckland", Country: "New Zealand", Category: "Oceania", Coordinates: [2]float64{-36.8485, 174.7633}, Color: tcell.ColorKhaki},
	{Name: "Sydney", Timezone: "Australia/Sydney", Country: "Australia", Category: "Oceania", Coordinates: [2]float64{-33.8688, 151.2093}, Color: tcell.ColorYellow},
}

// City aliases for common abbreviations
var cityAliases = map[string]string{
	"nyc":       "New York",
	"ny":        "New York",
	"lax":       "Los Angeles",
	"la":        "Los Angeles",
	"sf":        "San Francisco",
	"london":    "London",
	"lon":       "London",
	"paris":     "Paris",
	"par":       "Paris",
	"tokyo":     "Tokyo",
	"tyo":       "Tokyo",
	"singapore": "Singapore",
	"sin":       "Singapore",
	"dubai":     "Dubai",
	"dxb":       "Dubai",
	"sydney":    "Sydney",
	"syd":       "Sydney",
	"hong kong": "Hong Kong",
	"hk":        "Hong Kong",
	"mumbai":    "Mumbai",
	"bom":       "Mumbai",
	"shanghai":  "Shanghai",
	"pvg":       "Shanghai",
	"toronto":   "Toronto",
	"yvr":       "Vancouver",
	"mex":       "Mexico City",
	"sao paulo": "Sao Paulo",
	"gru":       "Sao Paulo",
	"buenos":    "Buenos Aires",
	"bue":       "Buenos Aires",
}

// GetCityByName looks up a city by name, handling aliases.
func GetCityByName(name string) *City {
	// First check aliases
	if canonical, ok := cityAliases[name]; ok {
		name = canonical
	}

	// Case-insensitive search
	lowerName := toLower(name)
	for i := range AllCities {
		if toLower(AllCities[i].Name) == lowerName {
			return &AllCities[i]
		}
	}
	return nil
}

// GetCitiesByCategory returns all cities in a specific category.
func GetCitiesByCategory(category string) []City {
	var result []City
	for _, city := range AllCities {
		if city.Category == category {
			result = append(result, city)
		}
	}
	return result
}

// GetCitiesByCategories returns cities from multiple categories.
func GetCitiesByCategories(categories []string) []City {
	var result []City
	for _, category := range categories {
		result = append(result, GetCitiesByCategory(category)...)
	}
	return result
}

// GetAllCategories returns the list of unique categories.
func GetAllCategories() []string {
	seen := make(map[string]bool)
	var categories []string
	for _, city := range AllCities {
		if !seen[city.Category] {
			seen[city.Category] = true
			categories = append(categories, city.Category)
		}
	}
	return categories
}

// toLower is a simple lowercase helper
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

// CategorizedCities returns a map of category to cities.
func CategorizedCities() map[string][]City {
	result := make(map[string][]City)
	for _, city := range AllCities {
		result[city.Category] = append(result[city.Category], city)
	}
	return result
}

// PrintCitiesList prints all cities grouped by category.
func PrintCitiesList() {
	fmt.Println("Available cities:")
	categories := GetAllCategories()
	for _, category := range categories {
		fmt.Printf("\n[%s]:\n", category)
		cities := GetCitiesByCategory(category)
		for _, city := range cities {
			fmt.Printf("  - %s (%s, %s)\n", city.Name, city.Timezone, city.Country)
		}
	}
	fmt.Println("\nAvailable presets:")
	for name := range presets {
		fmt.Printf("  - %s\n", name)
	}
}
