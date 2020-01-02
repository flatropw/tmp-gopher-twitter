package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func Response(w http.ResponseWriter, response map[string]interface{}) {
	w.Header().Add("Content-Type", "application-json")
	res, err := json.Marshal(&response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(res))
}
