package data

import (
	"unicode/utf8"

	"github.com/05blue04/Poneglyph/internal/validator"
)

var validRaces = map[string]struct{}{
	"human":            {},
	"fishman":          {},
	"merman":           {},
	"giant":            {},
	"dwarf":            {},
	"mink":             {},
	"lunarian":         {},
	"buccaneer":        {},
	"long arm tribe":   {},
	"long leg tribe":   {},
	"snake neck tribe": {},
	"three-eye tribe":  {},
	"snakeneck tribe":  {},
	"longarm tribe":    {},
	"longleg tribe":    {},
	"tontatta":         {},
	"kuja":             {},
	"skypiean":         {},
	"shandian":         {},
	"birkan":           {},
	"cyborg":           {},
	"zombie":           {},
	"artificial human": {},
	"reindeer":         {}, // For Chopper
	"skeleton":         {}, // For Brook
}

var validTypes = map[string]struct{}{
	"zoan":      {},
	"paramecia": {},
	"logia":     {},
}

func validateName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(len(name) < 300, "name", "must not be more than 300 bytes long")
	v.Check(utf8.ValidString(name), "name", "must be valid UTF-8")
}

func validateDescription(v *validator.Validator, description string) {
	v.Check(description != "", "description", "must be provided")
	v.Check(len(description) >= 10, "description", "must be at least 10 characters long")
	v.Check(len(description) <= 2000, "description", "must not be more than 2000 characters long")
	v.Check(utf8.ValidString(description), "description", "must be valid UTF-8")
}

func validateEpisode(v *validator.Validator, episode int) {
	v.Check(episode != 0, "episode", "must be provided")
	v.Check(episode <= 1200, "episode", "must not be greater than 1200")
	v.Check(episode > 0, "episode", "must not be negative")
}

func validateTimeSkip(v *validator.Validator, timeSkip string) {
	v.Check(timeSkip != "", "time_skip", "must be provided")
	v.Check(isValidTimeSkip(timeSkip), "time_skip", "must be either pre or post")
}

func IsValidRace(race string) bool {
	if race == "" {
		return false
	}

	_, exists := validRaces[race]
	return exists
}

func GetValidRaces() []string {
	races := make([]string, 0, len(validRaces))
	for race := range validRaces {
		races = append(races, race)
	}
	return races
}

func isValidTimeSkip(timeSkip string) bool {
	if timeSkip == "pre" {
		return true
	}

	if timeSkip == "post" {
		return true
	}

	return false
}

func IsValidType(devilFruitType string) bool {
	if devilFruitType == "" {
		return false
	}

	_, exists := validTypes[devilFruitType]

	return exists
}
