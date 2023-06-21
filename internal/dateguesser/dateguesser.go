package dateguesser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

//go:generate $MOCKGEN --source=dateguesser.go --destination=../../mocks/dateguesser/dateguesser_mock.go

// Global regex to match `-1`, `-2`, ...
var SubtractDayRegex = regexp.MustCompile(`^\-[0-9]+$`)

// IClock is a helper interface for useful time functions.
type IClock interface {
	Now() time.Time
}

// Clock implements IClock
type Clock struct{}

var _ IClock = &Clock{}

// Now implements IClock.
func (*Clock) Now() time.Time { return time.Now() }

// IDateGuesser proves an interface for an engine that guesses the date the user
// wants from it's text input.
type IDateGuesser interface {
	// Guess returns a tuple of (guess, success).
	Guess(userInput string) (guess time.Time, success bool)
}

// DateGuesser implemnets IDateGuesser.
type DateGuesser struct {
	Clock IClock
}

var _ IDateGuesser = &DateGuesser{}

// New returns a new instance of a DateGuesser
func New() (*DateGuesser, error) { return &DateGuesser{&Clock{}}, nil }

// Guess implements IDateGuesser
func (dg *DateGuesser) Guess(userInput string) (guess time.Time, success bool) {
	now := dg.Clock.Now()

	// Empty user input - use today
	if userInput == "" {
		return now, true
	}

	// User entered full date - we are happy
	if date, err := time.Parse("2006-01-02", userInput); err == nil {
		return date, true
	}

	// User entered partial day (without month/year)
	year := fmt.Sprint(now.Year())
	month := fmt.Sprint(int(now.Month()))
	if len(month) == 1 {
		month = "0" + month
	}
	if date, err := time.Parse("2006-01-02", year+"-"+month+"-"+userInput); err == nil {
		return date, true
	}

	// User entered partial month-day (without year)
	if date, err := time.Parse("2006-01-02", year+"-"+userInput); err == nil {
		return date, true
	}

	// User entered -x, where x is a number of days
	if match := SubtractDayRegex.MatchString(userInput); match {
		numOfDays, _ := strconv.Atoi(userInput[1:])
		return now.AddDate(0, 0, -numOfDays), true
	}

	return time.Time{}, false
}
