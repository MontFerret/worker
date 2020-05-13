package worker

type (
	// Query is the FQL-script.
	Query struct {
		Text   string                 `json:"text"`
		Params map[string]interface{} `json:"params"`
	}

	// Result is the result of Query.
	Result struct {
		Raw []byte
	}
)
