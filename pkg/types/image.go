package types

// Image comment
type Image struct {
	Payload   string   `json:"payload"`
	Encodings []string `json:"encodings"`
	Theme     string   `json:"theme"`
	Hash      string   `json:"hash"`
}
