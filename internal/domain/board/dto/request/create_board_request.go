package request

type CreateBoardRequest struct {
	Name       string  `json:"name"`
	Color      *string `json:"color"`
	Background *string `json:"background"`
}
