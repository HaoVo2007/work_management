package comment

type CommentService interface {}

type commentService struct {
	CommentRepository CommentRepository
}

func NewCommentService(commentRepository CommentRepository) CommentService {
	return &commentService{
		CommentRepository: commentRepository,
	}
}
