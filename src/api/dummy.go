package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type dummyData struct {
	Foo string `json:"foo"`
}

func GetDummyData(w http.ResponseWriter, r *http.Request) {
	body := []dummyData{{"bar"}, {"baz"}, {"qux"}}
	render.JSON(w, r, body)
}
