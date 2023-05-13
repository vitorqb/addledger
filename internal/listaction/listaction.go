package listaction

// A ListAction represents a possible action for a List input.
type ListAction string

const (
	NEXT ListAction = "NEXT"
	PREV ListAction = "PREV"
	NONE ListAction = "NONE"
)
