package hledger

// Client is the default implementation for IClient.
type Client struct {
	executable string
	ledgerFile string
}

// JSONTransaction represents a transaction in JSON
type JSONTransaction struct {
	Date        string        `json:"tdate"`
	Description string        `json:"tdescription"`
	Postings    []JSONPosting `json:"tpostings"`
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
