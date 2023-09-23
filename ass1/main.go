package main

import (
	"ass1/employee"
	"ass1/manager"
	"fmt"
)

func main() {

	m1 := manager.Manager{}
	m1.SetPosition("Manager")
	m1.SetSalary(60000.0)
	m1.SetAddress("123 Main St, City")

	// developer := &employee.Developer{}
	// developer.SetPosition("Developer")
	// developer.SetSalary(70000.0)
	// developer.SetAddress("456 Elm St, City")

	// // Print employee details
	printEmployeeDetails(&m1)
	// printEmployeeDetails(developer)
}

func printEmployeeDetails(e employee.Employee) {
	fmt.Printf("Position: %s\nSalary: %.2f\nAddress: %s\n", e.GetPosition(), e.GetSalary(), e.GetAddress())
}
