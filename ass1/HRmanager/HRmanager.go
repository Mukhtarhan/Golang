package hrHRmanager

type HRmanager struct {
	position string
	salary   float64
	address  string
}

func (m *HRmanager) GetPosition() string {
	return m.position
}

func (m *HRmanager) SetPosition(pos string) {
	m.position = pos
}

func (m *HRmanager) GetSalary() float64 {
	return m.salary
}

func (m *HRmanager) SetSalary(sal float64) {
	m.salary = sal
}

func (m *HRmanager) GetAddress() string {
	return m.address
}

func (m *HRmanager) SetAddress(addr string) {
	m.address = addr
}
