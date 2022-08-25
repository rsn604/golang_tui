package listdb
type ListItem struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Field01  string `json:"field01"`
	Field02  string `json:"field02"`
	Note     string `json:"note"`
}
