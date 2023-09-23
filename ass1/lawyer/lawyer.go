package Lawyer

type Lawyer struct {
	position string
	salary   float64
	address  string
}

func (m *Lawyer) GetPosition() string {
	return m.position
}

func (m *Lawyer) SetPosition(pos string) {
	m.position = pos
}

func (m *Lawyer) GetSalary() float64 {
	return m.salary
}

func (m *Lawyer) SetSalary(sal float64) {
	m.salary = sal
}

func (m *Lawyer) GetAddress() string {
	return m.address
}

func (m *Lawyer) SetAddress(addr string) {
	m.address = addr
}
