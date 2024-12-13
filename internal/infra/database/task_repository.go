package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"
const pageSize = 2

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Title       string            `db:"title"`
	Description string            `db:"description"`
	Status      domain.TaskStatus `db:"status"`
	Date        time.Time         `db:"date"`
	CreatedDate time.Time         `db:"created_date,omitempty"`
	UpdatedDate time.Time         `db:"updated_date,omitempty"`
	DeletedDate *time.Time        `db:"deleted_date,omitempty"`
}

type TaskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(sess db.Session) TaskRepository {
	return TaskRepository{
		coll: sess.Collection(TasksTableName),
		sess: sess,
	}
}

func (r TaskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()

	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(tsk), nil
}

func (r TaskRepository) Find(id uint64) (domain.Task, error) {
	var tsk task

	err := r.coll.Find(db.Cond{"id": id}).One(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(tsk), nil
}

func (r TaskRepository) FindByUser(uid uint64) ([]domain.Task, error) {
	var tasks []task

	err := r.coll.Find(db.Cond{"user_id": uid}).All(&tasks)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(tasks), nil
}

func (r TaskRepository) Update(t domain.Task, id, uid uint64) (domain.Task, error) {
	var taskToUpdate domain.Task

	err := r.coll.Find(db.Cond{"id": id, "user_id": uid}).One(&taskToUpdate)
	if err != nil {
		return domain.Task{UserId: uid}, err
	}

	taskToReturn := r.mapDomainToModel(t)
	taskToReturn.Id = id
	taskToReturn.UserId = uid
	taskToReturn.UpdatedDate = time.Now()
	taskToReturn.CreatedDate = taskToUpdate.CreatedDate
	taskToReturn.DeletedDate = taskToUpdate.DeletedDate

	err = r.coll.UpdateReturning(&taskToReturn)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(taskToReturn), nil
}

func (r TaskRepository) Delete(id uint64) (domain.Task, error) {
	var tsk task

	err := r.coll.Find(db.Cond{"id": id}).One(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	now := time.Now()
	tsk.DeletedDate = &now

	err = r.coll.UpdateReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(tsk), nil
}

func (r TaskRepository) mapDomainToModel(d domain.Task) task {
	return task{
		Id:          d.Id,
		UserId:      d.UserId,
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
		Date:        d.Date,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r TaskRepository) mapModelToDomain(m task) domain.Task {
	return domain.Task{
		Id:          m.Id,
		UserId:      m.UserId,
		Title:       m.Title,
		Description: m.Description,
		Status:      m.Status,
		Date:        m.Date,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r TaskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	var tasks []domain.Task
	for _, t := range ts {
		tasks = append(tasks, r.mapModelToDomain(t))
	}
	return tasks
}
