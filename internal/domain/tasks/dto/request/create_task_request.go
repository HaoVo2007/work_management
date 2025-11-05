package request

type CreateTaskRequest struct {
	ColumnID    string `json:"column_id"`
	Title       string `json:"title"`
}