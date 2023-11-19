package main

import (
	utils "github.com/niewolinsky/tw_employee_data_processor/utils"

	"log/slog"
	"net/http"
)

func (app *application) hdlGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	err := utils.WriteJSON(w, http.StatusOK, utils.Wrap{"status": "Status OK, autodeployed 7"}, nil)
	if err != nil {
		slog.Error("Unable to send healthcheckHandler response", err)
	}
}
