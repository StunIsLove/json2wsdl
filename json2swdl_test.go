package main

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandleSOAPRequest(t *testing.T) {
	// Создаем тестовый JSON запрос
	jsonData := []byte(`{
		"bik": "1234567890",
		"fromDate": "07-11-2018",
		"toDate": "07-18-2019",
		"withDeleted": true
	}`)

	// Создаем тестовый запрос
	req, err := http.NewRequest("POST", "http://localhost:8080/api/wsdl", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Записываем тело запроса
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder (реализация http.ResponseWriter) для записи ответа
	rr := httptest.NewRecorder()

	// Создаем маршрутизатор и регистрируем наш обработчик
	r := mux.NewRouter()
	r.HandleFunc("/api/wsdl", handleSOAPRequest).Methods(http.MethodPost)

	// Выполняем запрос
	r.ServeHTTP(rr, req)

	// Проверяем код статуса
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, but got %v", rr.Code)
	}

	// Парсим полученный SOAP-ответ в структуру SOAPResponse
	var soapResp SOAPResponse
	err = xml.Unmarshal(rr.Body.Bytes(), &soapResp)
	if err != nil {
		t.Errorf("Error parsing SOAP response: %v", err)
	}

	// Проверяем, что значения соответствуют ожидаемым
	expectedBik := "1234567890"
	expectedFromDate := "07-11-2018"
	expectedToDate := "07-18-2019"
	expectedWithDeleted := true

	print(soapResp.Bik)
	print(soapResp.FromDate)
	print(soapResp.ToDate)
	print(soapResp.WithDeleted)

	if soapResp.Bik != expectedBik {
		t.Errorf("Expected bik to be %s, but got %s", expectedBik, soapResp.Bik)
	}

	if soapResp.FromDate != expectedFromDate {
		t.Errorf("Expected fromDate to be %s, but got %s", expectedFromDate, soapResp.FromDate)
	}

	if soapResp.ToDate != expectedToDate {
		t.Errorf("Expected toDate to be %s, but got %s", expectedToDate, soapResp.ToDate)
	}

	if soapResp.WithDeleted != expectedWithDeleted {
		t.Errorf("Expected withDeleted to be %t, but got %t", expectedWithDeleted, soapResp.WithDeleted)
	}
}

func TestInvalidJSONData(t *testing.T) {
	// Создаем тестовый некорректный JSON запрос
	jsonData := []byte(`{
		"bik": "1234567890111",
		"fromDate": "32-11-2018",
		"toDate": "07-18-2019"
	}`)

	// Создаем тестовый запрос
	req, err := http.NewRequest("POST", "http://localhost:8080/api/wsdl", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Записываем тело запроса
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder (реализация http.ResponseWriter) для записи ответа
	rr := httptest.NewRecorder()

	// Создаем маршрутизатор и регистрируем наш обработчик
	r := mux.NewRouter()
	r.HandleFunc("/api/wsdl", handleSOAPRequest).Methods(http.MethodPost)

	// Выполняем запрос
	r.ServeHTTP(rr, req)

	// Проверяем код статуса (должен быть 400 Bad Request)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, but got %v", rr.Code)
	}
}

func TestMain(m *testing.M) {
	// Выполнение тестов
	code := m.Run()

	os.Exit(code)
}
