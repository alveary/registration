package announce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// NewService provides a method to attach a new Service to the overseer stack
func NewService(serviceName string, serviceRoot string, aliveResource string) {
	json, _ := json.Marshal(struct {
		serviceName   string
		serviceRoot   string
		aliveResource string
	}{
		serviceName,
		serviceRoot,
		aliveResource,
	})

	overseerRoot := os.Getenv("OVERSEER_ROOT")

	if overseerRoot == "" {
		fmt.Println("OVERSEER_ROOT is not set:")
		return
	}

	http.Post(overseerRoot, "application/json", bytes.NewBuffer(json))
}
