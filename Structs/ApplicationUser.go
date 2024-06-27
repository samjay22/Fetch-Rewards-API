package Structs

// Application User struct, scans are in the form of id in pattern "^\\S+$"
type User struct {
	Id     string   `json:"id"`
	Points int64    `json:"points"`
	Scans  []string `json:"scans"`
}
