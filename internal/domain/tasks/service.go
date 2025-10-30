package tasks

type TaskService interface {

}

type taskService struct {
	TaskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) TaskService {
	return &taskService{
		TaskRepository: taskRepository,
	}
}
