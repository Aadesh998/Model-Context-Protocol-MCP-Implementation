package tools

import (
	"main/utils"
)

type GetEmpLeavesInput struct {
	EmployeeID string `json:"employee_id" jsonschema_description:"The ID of the employee."`
}

var GetEmpLeavesInputSchema = utils.GenerateSchema[GetEmpLeavesInput]()

type GetEmpLeavesResponse struct {
	EmployeeID string `json:"employee_id"`
	LeaveType  string `json:"leave_type"`
	LeaveDate  string `json:"leave_date"`
	Approved   bool   `json:"approved"`
}

func GetEmpLeaves(employeeID string) GetEmpLeavesResponse {
	if employeeID == "1234" {
		return GetEmpLeavesResponse{
			EmployeeID: "1234",
			LeaveType:  "Sick Leave",
			LeaveDate:  "2025-06-20",
			Approved:   true,
		}
	}

	return GetEmpLeavesResponse{
		EmployeeID: employeeID,
		LeaveType:  "Unknown",
		LeaveDate:  "N/A",
		Approved:   false,
	}
}
