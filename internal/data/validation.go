package data

import (
	"unicode/utf8"

	"github.com/05blue04/Poneglyph/internal/validator"
)

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
