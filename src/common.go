package src

import (
	guuid "github.com/satori/go.uuid"
	"net/http"
)

// OptionsAnswer create options answer for browser
func OptionsAnswer(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
}

// FillAnswerHeader add some important headers to answer
func FillAnswerHeader(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
}

// GetHash we use it to get hasj=h for todo command
func GetHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}
