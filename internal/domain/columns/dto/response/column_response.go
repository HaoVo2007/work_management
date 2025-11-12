package response

type ColumnResponse struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Position int             `json:"position"`
	Tasks    []*TaskResponse `json:"tasks"`
}

type TaskResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	ImageCover *string `json:"image_cover"`
}
