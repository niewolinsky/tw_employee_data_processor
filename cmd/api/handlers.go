package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	customerimporter "github.com/niewolinsky/customerimporter"
	utils "github.com/niewolinsky/tw_employee_data_processor/utils"
	pb "github.com/niewolinsky/tw_employee_data_service/grpc/employee"
)

func (app *application) hdlGetUniqueEmails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.ListEmployeesRequest{
		Limit:  10000,
		Offset: 0,
	}

	// Use the gRPC client from the application struct
	employees, err := app.grpcEmployeeClient.ListEmployees(ctx, req)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	var domainProviders []customerimporter.DomainProvider
	for _, employee := range employees.GetEmployees() {
		domainProviders = append(domainProviders, EmployeeResponseWrapper{employee})
	}

	uniqueDomainsSorted := customerimporter.CountDomainsConcurrent(domainProviders)

	fmt.Printf("%v", uniqueDomainsSorted)

	// Write the response as JSON
	err = utils.WriteJSON(w, http.StatusOK, utils.Wrap{"uniqueDomainsSorted": uniqueDomainsSorted}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}
}
