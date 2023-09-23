// employee/manager.go
package manager

type Manager struct {
	position string
	salary   float64
	address  string
}

func (m *Manager) GetPosition() string {
	return m.position
}

func (m *Manager) SetPosition(pos string) {
	m.position = pos
}

func (m *Manager) GetSalary() float64 {
	return m.salary
}

func (m *Manager) SetSalary(sal float64) {
	m.salary = sal
}

func (m *Manager) GetAddress() string {
	return m.address
}

func (m *Manager) SetAddress(addr string) {
	m.address = addr
}
