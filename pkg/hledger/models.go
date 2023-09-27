package hledger

// JSONTag represents a tag in JSON
type JSONTag struct {
	Name  string
	Value string
}

// JSONTransaction represents a transaction in JSON
type JSONTransaction struct {
	Date        string        `json:"tdate"`
	Description string        `json:"tdescription"`
	Comment     string        `json:"tcomment"`
	Postings    []JSONPosting `json:"tpostings"`
	Tags        []JSONTag     `json:"ttags"`
}

// JSONPosting represents a posting in JSON
type JSONPosting struct {
	Account string        `json:"paccount"`
	Ammount []JSONAmmount `json:"pamount"`
}

// JSONAmmount represents an ammount in JSON
type JSONAmmount struct {
	Commodity string       `json:"acommodity"`
	Quantity  JSONQuantity `json:"aquantity"`
}

// JSONQuantity represents a quantity in JSON
type JSONQuantity struct {
	DecimalMantissa int64 `json:"decimalMantissa"`
	DecimalPlaces   int32 `json:"decimalPlaces"`
}
