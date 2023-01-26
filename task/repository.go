package task

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Store(inputTask InputTask) (Task, error)
	FetchAll() ([]Task, error)
	FetchById(id int) (Task, error)
	Update(id int, inputTask InputTask) (Task, error)
	Delete(id int) error
}

type repositry struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	repo := repositry{db}

	return &repo
}

func (r *repositry) Store(inputTask InputTask) (Task, error) {
	task := Task{
		Name:        inputTask.Name,
		Description: inputTask.Description,
	}
	err := r.db.Create(&task).Error
	if err != nil {
		return task, err
	}

	return task, nil
}

func (r *repositry) FetchAll() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *repositry) FetchById(id int) (Task, error) {
	task := Task{ID: id}
	err := r.db.First(&task).Error

	if err != nil {
		return task, err
	}

	return task, nil
}

func (r *repositry) Update(id int, inputTask InputTask) (Task, error) {
	task, err := r.FetchById(id)
	if err != nil {
		return task, err
	}

	task.Name = inputTask.Name
	task.Description = inputTask.Description
	err = r.db.Save(&task).Error
	if err != nil {
		return task, err
	}

	return task, nil
}

func (r *repositry) Delete(id int) error {
	tx := r.db.Delete(Task{ID: id})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("id : %v not found", id)
	}

	return nil
}
