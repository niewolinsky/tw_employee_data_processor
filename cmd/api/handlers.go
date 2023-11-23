package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	customerimporter "github.com/niewolinsky/customerimporter"
	utils "github.com/niewolinsky/tw_employee_data_processor/utils"
	pb "github.com/niewolinsky/tw_employee_data_service/grpc/employee"
)

func (app *application) hdlGetUniqueEmails(w http.ResponseWriter, r *http.Request) {
	cachedResponse, err := app.redisClient.Get(r.Context(), "uniqueDomainsSorted").Result()
	if err != nil {
		switch {
		case (err.Error() == "redis: nil"):
			fmt.Println("empty cache")
		default:
			slog.Error("cache error", "MESSAGE", err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(cachedResponse))
		if err != nil {
			utils.ServerErrorResponse(w, r, err)
		}
		return
	}

	req := &pb.ListEmployeesRequest{
		Limit:  10000,
		Offset: 0,
	}

	employees, err := app.grpcEmployeeClient.ListEmployees(r.Context(), req)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	var domainProviders []customerimporter.DomainProvider
	for _, employee := range employees.GetEmployees() {
		domainProviders = append(domainProviders, EmployeeResponseWrapper{employee})
	}

	uniqueDomainsSorted := customerimporter.CountDomainsConcurrent(domainProviders)

	jsonData, err := utils.WriteJSONCache(w, http.StatusOK, utils.Wrap{"uniqueDomainsSorted": uniqueDomainsSorted}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}

	err = app.redisClient.Set(r.Context(), "uniqueDomainsSorted", jsonData, time.Hour*1).Err()
	if err != nil {
		slog.Error("failed caching response", "MESSAGE", err)
	}
}

func (app *application) hdlPatchUniqueEmails(w http.ResponseWriter, r *http.Request) {
	cacheKey := "uniqueDomainsSorted"

	err := app.redisClient.Del(r.Context(), cacheKey).Err()
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Wrap{"message": "Employee unique domains cache data cleared."}, nil)
}
