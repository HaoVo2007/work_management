package request

type UpdateColumnRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}