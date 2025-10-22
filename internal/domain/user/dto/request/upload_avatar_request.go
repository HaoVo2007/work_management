package request
type UploadAvatarRequest struct {
	Avatar string `json:"avatar" form:"avatar"`
}
