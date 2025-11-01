package request

type UpdateBoardRequest struct {
	Name       string  `json:"name"`
	Color      *string `json:"color"`
	Icon       *string `json:"icon"`
	Background *string `json:"background"`
}
