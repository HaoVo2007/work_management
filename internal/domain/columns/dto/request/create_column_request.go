package request

type CreateColumnRequest struct {
	BoardID  string `json:"board_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
}
