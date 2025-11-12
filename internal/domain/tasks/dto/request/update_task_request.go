package request

import "mime/multipart"

type UpdateTaskRequest struct {
	Title       *string               `form:"title"`
	Description *string               `form:"description"`
	CoverPhoto  *multipart.FileHeader `form:"cover_photo"`
	Assignee    *string               `form:"assignee"`
	StartDate   *string               `form:"start_date"`
	EndDate     *string               `form:"end_date"`
	Priority    *int64                `form:"priority"`
}
