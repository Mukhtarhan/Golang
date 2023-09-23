package developer

type Developer struct {
	position string
	salary   float64
	address  string
}

func (m *Developer) GetPosition() string {
	return m.position
}

func (m *Developer) SetPosition(pos string) {
	m.position = pos
}

func (m *Developer) GetSalary() float64 {
	return m.salary
}

func (m *Developer) SetSalary(sal float64) {
	m.salary = sal
}

func (m *Developer) GetAddress() string {
	return m.address
}

func (m *Developer) SetAddress(addr string) {
	m.address = addr
}
