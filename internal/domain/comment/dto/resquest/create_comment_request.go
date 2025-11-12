package resquest

import "mime/multipart"

type CreateCommentRequest struct {
	Content *string               `form:"content"`
	File    *multipart.FileHeader `form:"file"`
}
