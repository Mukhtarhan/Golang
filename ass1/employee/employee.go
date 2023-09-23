package employee

type Employee interface {
	GetPosition() string
	SetPosition(string)
	GetSalary() float64
	SetSalary(float64)
	GetAddress() string
	SetAddress(string)
}
