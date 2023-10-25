package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type InputData struct {
	Bik         string `json:"bik"`
	FromDate    string `json:"fromDate"`
	ToDate      string `json:"toDate"`
	WithDeleted bool   `json:"withDeleted"`
}

type SOAPResponse struct {
	XMLName     xml.Name `xml:"soapResponse"`
	Bik         string   `xml:"bik"`
	FromDate    string   `xml:"fromDate"`
	ToDate      string   `xml:"toDate"`
	WithDeleted bool     `xml:"withDeleted"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/wsdl", handleSOAPRequest).Methods(http.MethodPost)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func handleSOAPRequest(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON из тела запроса
	var inputData InputData
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if inputData.Bik == "" || (len(inputData.Bik) != 10 && len(inputData.Bik) != 12) {
		http.Error(w, "Invalid 'bik' field", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("01-02-2006", inputData.FromDate)
	if err != nil {
		http.Error(w, "Invalid 'fromDate' field", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("01-02-2006", inputData.ToDate)
	if err != nil {
		http.Error(w, "Invalid 'toDate' field", http.StatusBadRequest)
		return
	}

	// Генерируем SOAP ответ
	soapResponse := SOAPResponse{
		Bik:         inputData.Bik,
		FromDate:    inputData.FromDate,
		ToDate:      inputData.ToDate,
		WithDeleted: inputData.WithDeleted,
	}

	// Преобразуем SOAP структуру в XML
	xmlResponse, err := xml.MarshalIndent(soapResponse, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Устанавливаем необходимые заголовки для ответа
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Возвращаем XML как ответ
	w.Write(xmlResponse)
}
