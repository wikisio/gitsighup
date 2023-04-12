package action

import (
	"net/http"
	"testing"
)

func TestUrl(t *testing.T) {
	url := "http://localhost:8080/" + namespace + "/" + service + "/" + filename
	resp, err := http.Get(url)
}
