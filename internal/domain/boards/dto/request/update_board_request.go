package request

type UpdateBoardRequest struct {
	Name       string  `json:"name"`
	Color      *string `json:"color"`
	Background *string `json:"background"`
}
