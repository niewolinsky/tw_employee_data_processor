package main

import (
	"strings"

	pb "github.com/niewolinsky/tw_employee_data_service/grpc/employee"
)

type EmployeeResponseWrapper struct {
	*pb.EmployeeResponse
}

func (e EmployeeResponseWrapper) GetDomain() string {
	parts := strings.Split(string(e.GetEmail()), "@")
	return parts[1]
}
