package request

type CreateBoardRequest struct {
	Name       string  `json:"name"`
	Icon       *string `json:"icon"`
	Color      *string `json:"color"`
	Background *string `json:"background"`
}
