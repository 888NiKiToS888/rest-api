package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// обработчик всех задач
/*Обработчик должен вернуть все задачи, которые хранятся в мапе.
Конечная точка /tasks.
Метод GET.
При успешном запросе сервер должен вернуть статус 200 OK.
При ошибке сервер должен вернуть статус 500 Internal Server Error.
*/

func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	_, err = w.Write(resp) //проверка на ошибку
	if err != nil {
		fmt.Printf("ошибка при записи в тело ответа: %s\n", err.Error()) // логируем ошибку
	}
}

// обработчик новой задачи
/*Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
Конечная точка /tasks.
Метод POST.
При успешном запросе сервер должен вернуть статус 201 Created.
При ошибке сервер должен вернуть статус 400 Bad Request.
*/

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf []byte
	var err error

	buf, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf, &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := tasks[task.ID]; ok {
		http.Error(w, fmt.Sprintf("Задача с таким ID %s уже существует", task.ID), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

// обработчик задачи по ID
/*Обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе.
В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.
Конечная точка /tasks/{id}.
Метод GET.
При успешном выполнении запроса сервер должен вернуть статус 200 OK.
В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
*/

func getTaskId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]

	if !ok {
		http.Error(w, "не найден", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Printf("ошибка при записи ответа: %s\n", err.Error())
	}
}

// обработчик удаления задачи по ID
/*Обработчик должен удалить задачу из мапы по её ID. Здесь так же нужно сначала проверить, есть ли задача с таким ID в мапе, если нет вернуть соответствующий статус.
Конечная точка /tasks/{id}.
Метод DELETE.
При успешном выполнении запроса сервер должен вернуть статус 200 OK.
В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
*/

func deleteTaskId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]

	if !ok {
		http.Error(w, "не найден", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики

	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTaskId)
	r.Delete("/tasks/{id}", deleteTaskId)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
