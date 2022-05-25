package middleware

import (
	"encoding/json"
	"go-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetAllTask(t *testing.T) {
	task1 := models.ToDoList{primitive.NewObjectID(), "Task1", false}
	task2 := models.ToDoList{primitive.NewObjectID(), "Task2", false}
	SetTaskRegistryForTests(testTaskRegistry{
		task1.ID.Hex(): &task1,
		task2.ID.Hex(): &task2,
	})

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/task", nil)
	GetAllTask(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	}
	actual := unmarshal(res)
	expected := []models.ToDoList{task1, task2}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got body %v but wanted %v", actual, expected)
	}
}

func unmarshal(res *httptest.ResponseRecorder) []models.ToDoList {
	var actual []models.ToDoList
	json.Unmarshal(res.Body.Bytes(), &actual)
	return actual
}

type testTaskRegistry map[string]*models.ToDoList

func (tr testTaskRegistry) GetAllTask() []models.ToDoList {
	entries := make([]models.ToDoList, 0, len(tr))
	for _, e := range tr {
		entries = append(entries, *e)
	}
	return entries
}

func (tr testTaskRegistry) InsertOneTask(task models.ToDoList) {
	tr[task.ID.Hex()] = &task
}

func (tr testTaskRegistry) TaskComplete(task string) {
	tr[task].Status = true
}

func (tr testTaskRegistry) UndoTask(task string) {
	tr[task].Status = false
}

func (tr testTaskRegistry) DeleteOneTask(task string) {
	delete(tr, task)
}
