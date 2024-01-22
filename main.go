package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"github.com/webbben/code-exec-microservice/execute"
)

type ExecRequest struct {
	Code string `json:"code"`
	Lang string `json:"lang"`
}
type ExecResponse struct {
	Output string `json:"output"`
	Error  bool   `json:"error"`
}

var supportedLangs = []string{"python", "go", "bash"}

func main() {
	r := mux.NewRouter()

	// TODO implement authentication to limit who can use this API
	r.HandleFunc("/", handleExecRequest).Methods("POST")

	fmt.Println("== code execution service! ==")
	fmt.Println("Server listening on localhost:8081")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleExecRequest(w http.ResponseWriter, r *http.Request) {
	var req ExecRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if !slices.Contains(supportedLangs, req.Lang) {
		http.Error(w, fmt.Sprintf("Language %s not supported", req.Lang), http.StatusBadRequest)
		return
	}
	log.Printf("received %s code execution request", req.Lang)
	output, err := execute.ExecuteCode(req.Lang, req.Code, false)
	var res ExecResponse
	if err != nil {
		res = ExecResponse{
			Output: err.Error(),
			Error:  true,
		}
	} else {
		res = ExecResponse{
			Output: output,
			Error:  false,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
