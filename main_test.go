package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	w := httptest.NewRecorder() //создает responsewriter для тестов

	handleRoot(w, nil) // запускаем функцию

	desiredCode := http.StatusOK

	if w.Code != desiredCode { // проверяем чтобы статус возвращаелся таким, что все работает
		t.Errorf("bad response code, expected %v but got %v\nbode: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Welcome to our homepage!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) { // также проверяем чтобы текст на сайте соответствовал тому, что мы отправляли
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleGoodbye(t *testing.T) {
	w := httptest.NewRecorder()

	handleGoodbye(w, nil)

	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\nbode: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Goodbye!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected %q", w.Body.String(), expectedMessage)
	}

}

func TestHandleHelloParameterized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?user=TestMan", nil) //create request with method get and URL with parameter

	w := httptest.NewRecorder()

	handleHelloParameterized(w, req) // put responsewriter and request

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected: %q",
			w.Body.String(), expectedMessage)
	}
}

func TestHandleHelloParameterizedNoParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello/", nil)

	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected: %q",
			w.Body.String(), expectedMessage)
	}
}

func TestHandleParameterizedWrongParameter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?foo=bar", nil)

	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected: %q",
			w.Body.String(), expectedMessage)
	}
}

func TestHandleUserResponsesHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/responses/Testman/hello/", nil)
	req.SetPathValue("user", "Testman") // set url variable "user" with value "Testman"

	w := httptest.NewRecorder()

	handleUserResponsesHello(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, Testman!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}
}

func TestHandleHelloHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello", nil)
	req.Header.Set("user", "Test Man") //set header

	w := httptest.NewRecorder()

	handleHelloHeader(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, Test Man!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}

}

func TestHandleHelloNoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello", nil)

	w := httptest.NewRecorder()

	handleHelloHeader(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("invalid username provided\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}
}

func TestHandleJSON(t *testing.T) {
	testRequest := UserData{
		Name: "Test Man",
	}
	marshaledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshaledRequestBody))

	w := httptest.NewRecorder()

	handleJSON(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, Test Man!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}
}

func TestHandleJSONEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/json", nil)

	w := httptest.NewRecorder()

	handleJSON(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("bad request body\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}
}

func TestHandleJSONEmptyNameFiled(t *testing.T) {
	testRequest := UserData{
		Name: "",
	}
	marshaledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshaledRequestBody))

	w := httptest.NewRecorder()

	handleJSON(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("invalid username provided\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return got: %q, expected %q", w.Body.Bytes(), expectedMessage)
	}
}
