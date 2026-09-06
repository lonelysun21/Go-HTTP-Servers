package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type UserData struct {
	Name string
}

func main() {
	mux := http.NewServeMux() //Наш маршрутизатор, в котором мы храним нашу маршруты (адрес и его функцию)

	mux.HandleFunc("/{$}", handleRoot) //Наши маршруты в которых хранятся адреса и их функции
	mux.HandleFunc("/goodbye", handleGoodbye)
	mux.HandleFunc("/hello/", handleHelloParameterized)
	mux.HandleFunc("/responses/{user}/hello/", handleUserResponsesHello)
	mux.HandleFunc("/user/hello/", handleHelloHeader)
	mux.HandleFunc("/json/", handleJSON)

	http.ListenAndServe(":8080", mux) // Нужен чтобы подключиться к серверу, запустить его как бесконечный цикл и принять маршрутизатор который отправиться на сервер
}

// handleRoot for homepage with message Welcome
func handleRoot(w http.ResponseWriter, r *http.Request) { // в каждом хэндлере у нас всегда 2 параметра, w and r,
	// в котором в w у нас храниться то что мы будем отправлять на страничку, а в r то что мы имеем и можем использовать со странички (body, code and etc)
	_, err := w.Write([]byte("Welcome to our homepage!\n")) // Функция w.write принимает слайс байтов и нужно чтобы писать текст на страничке.
	if err != nil {
		slog.Error("Error writing responce", "err", err)
		return
	}
}

// handleGoodbye for Goodbye message
func handleGoodbye(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Goodbye!\n"))
	if err != nil {
		slog.Error("Error writing responce", "err", err)
		return
	}
}

// handleHelloParameterized for taking parameter from URL and writing it on the website page
func handleHelloParameterized(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()    //taking parameters from URL
	userList := params["user"] //taking parameter with name "user"

	username := "User"
	if len(userList) > 0 { //if this user parameter and its value exist we use it, else use "User".
		username = userList[0]
	}

	handleHello(w, username)
}

// handle UserResponsesHello for working with url variables
func handleUserResponsesHello(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user") //taking url variable with name "user"

	handleHello(w, username)
}

// handleHelloHeader for working with headers
func handleHelloHeader(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("user") // taking value from header with key "user"
	if username == "" {
		http.Error(w, "invalid username provided", http.StatusBadRequest)
		return
	}

	handleHello(w, username)
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	byteData, err := io.ReadAll(r.Body)
	if err != nil || len(byteData) < 1 {
		slog.Error("error reading request body", "err", nil)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	var reqData UserData
	err = json.Unmarshal(byteData, &reqData)
	if err != nil {
		slog.Error("error unmarshalling request body", "err", err)
		http.Error(w, "error parsing request JSON", http.StatusBadRequest)
		return
	}

	if len(reqData.Name) == 0 {
		slog.Error("error empty field", "err", err)
		http.Error(w, "invalid username provided", http.StatusBadRequest)
		return
	}

	handleHello(w, reqData.Name)
}

// handleHello that write Hello and user nickname on the site
func handleHello(w http.ResponseWriter, username string) {
	var output bytes.Buffer //buffer save strings
	output.WriteString("Hello, ")
	output.WriteString(username)
	output.WriteString("!\n")

	_, err := w.Write(output.Bytes())
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}
