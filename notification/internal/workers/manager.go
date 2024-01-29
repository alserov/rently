package workers

type Manager interface {
	Add(worker Worker)
}

func NewManager() Manager {
	return &manager{}
}

type manager struct {
	workers []Worker
}

func (m manager) Add(worker Worker) {
	worker.Start()
	m.workers = append(m.workers, worker)
}
