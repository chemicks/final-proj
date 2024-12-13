package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save -> requests.Bind: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id

		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save -> c.taskService.Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Created(w, taskDto)
	}
}

func (c TaskController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("FindById: invalid %s parameter(only non-negative integers)", chi.URLParam(r, "id"))
			BadRequest(w, err)
			return
		}

		task, err := c.taskService.Find(id)
		if err != nil {
			log.Printf("TaskController -> FindById -> c.taskService.Find: %s", err)
			InternalServerError(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		if task.UserId != user.Id {
			err = errors.New("access denied")
			Forbidden(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) FindByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		tasks, err := c.taskService.FindByUser(user.Id)
		if err != nil {
			log.Printf("TaskController -> FindByUser -> c.taskService.FindByUser: %s", err)
			InternalServerError(w, err)
			return
		}

		var tasksDto resources.TasksDto
		tasksDto = tasksDto.DomainToDto(tasks)
		Success(w, tasksDto)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("Update: invalid %s parameter(only non-negative integers)", chi.URLParam(r, "id"))
			BadRequest(w, err)
			return
		}

		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Update -> requests.Bind: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)

		task, err = c.taskService.Update(task, id, user.Id)
		if task.UserId != user.Id {
			err = errors.New("access denied")
			Forbidden(w, err)
			return
		}
		if err != nil {
			log.Printf("TaskController -> Update -> c.taskService.Update: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("Delete: invalid %s parameter(only non-negative integers)", chi.URLParam(r, "id"))
			BadRequest(w, err)
			return
		}

		task, err := c.taskService.Delete(id)
		if err != nil {
			log.Printf("TaskController -> Delere -> c.taskService.Delere: %s", err)
			InternalServerError(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		if task.UserId != user.Id {
			err = errors.New("access denied")
			Forbidden(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}
