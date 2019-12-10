package src

import (
	"net/http"

	guuid "github.com/satori/go.uuid"
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

// FindOccurences find second occurences of a char
func FindOccurences(source, char string, occur int) int {
	count := 0
	for i := 0; i < len(source); i++ {
		if source[i:i+1] == char {
			count++
			if count == occur {
				return i
			}
		}
	}
	return len(source)
}

// GetRoute get main route
func GetRoute(inRoute string) string {
	ind := FindOccurences(inRoute, "/", 2)
	return inRoute[:ind]
}
