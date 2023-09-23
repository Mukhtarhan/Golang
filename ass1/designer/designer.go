package designer

type Designer struct {
	position string
	salary   float64
	address  string
}

func (m *Designer) GetPosition() string {
	return m.position
}

func (m *Designer) SetPosition(pos string) {
	m.position = pos
}

func (m *Designer) GetSalary() float64 {
	return m.salary
}

func (m *Designer) SetSalary(sal float64) {
	m.salary = sal
}

func (m *Designer) GetAddress() string {
	return m.address
}

func (m *Designer) SetAddress(addr string) {
	m.address = addr
}
