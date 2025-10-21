package user

type UploadAvatarRequest struct {
	Avatar string `json:"avatar" form:"avatar"`
}