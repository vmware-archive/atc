package api

import "github.com/concourse/atc"

// WorkerService declares interface for worker business functionality
type WorkerService interface {
	Land(worker atc.Worker) error
	Retire(worker atc.Worker) error
	Prune(worker atc.Worker) error
	Delete(worker atc.Worker) error
}
