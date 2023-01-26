package task_test

import (
	"fmt"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/MohammadBnei/iw-janvier-task/task"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db = initDB()

func initDB() *gorm.DB {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "user:password@tcp(127.0.0.1:3306)/go-course?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&task.Task{})

	return db
}

func cleanUp() {
	fmt.Println(db.Where("1 = 1").Delete(task.Task{}).RowsAffected)
}

func TestCreate(t *testing.T) {
	repo := task.NewRepository(db)
	defer cleanUp()

	task, err := repo.Store(task.InputTask{Name: "First Test", Description: "First Test description"})
	if err != nil {
		t.Error(err)
	}

	if task.Name != "First Test" || task.Description != "First Test description" {
		t.Error("Task not correctly saved")
	}
}

func TestFindAll(t *testing.T) {
	repo := task.NewRepository(db)
	defer cleanUp()

	for i := 0; i < 10; i++ {
		repo.Store(task.InputTask{
			Name: fmt.Sprintf("task %v", i),
		})
	}

	tasks, err := repo.FetchAll()
	if err != nil {
		t.Error(err)
	}

	if len(tasks) != 10 {
		t.Errorf("Not enought tasks, wanted 10 got %v", len(tasks))
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	for i, task := range tasks {
		if task.Name != fmt.Sprintf("task %v", i) {
			t.Errorf("Wrong task name, wanted %v got %v", fmt.Sprintf("task %v", i), task.Name)
		}
	}
}
