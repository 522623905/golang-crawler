package controller

import (
	"github.com/fatrbaby/cobweb/web/view"
	"net/http"
)

type HomeHandler struct {
	view view.HomeView
}

func (that HomeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := that.view.Render(response)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}
}

func NewHomeHandler(template string) HomeHandler {
	return HomeHandler{
		view: view.CreateHomeView(view.Load(template)),
	}
}
