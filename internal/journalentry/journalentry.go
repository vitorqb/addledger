package journalentry

import (
	"time"
)

type (
	JournalEntry struct {
		Date        time.Time
		Description string
	}
)
