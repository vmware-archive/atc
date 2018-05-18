package service

import (
	"time"

	"github.com/concourse/atc"
)

//go:generate counterfeiter . WorkerRepository

// WorkerRepository declares the interfaces for persistent storage of Worker information
type WorkerRepository interface {
	GetWorker(name string) (*atc.Worker, bool, error)
	SaveWorker(atcWorker atc.Worker, ttl time.Duration) (*atc.Worker, error)
	HeartbeatWorker(worker atc.Worker, ttl time.Duration) (*atc.Worker, error)
	Workers() ([]atc.Worker, error)
	VisibleWorkers([]string) ([]atc.Worker, error)
	Reload(*atc.Worker) (bool, error)
	Delete(*atc.Worker) error
}

// Worker struct that anchors the implementation of the WorkerService interface
type Worker struct {
	WorkerRepo WorkerRepository
}

// Land lands the worker
// Worker will be unavailable for scheduling while it's in landed state
func (workerService Worker) Land(worker atc.Worker) error {
	return nil
}

// Retire removes the worker from the pool permanently
func (workerService Worker) Retire(worker atc.Worker) error {
	return nil
}

// Prune will remove a worker from the worker pool after it has been found to be stalled
func (workerService Worker) Prune(worker atc.Worker) error {
	return nil
}

// Delete will remove the worker TODO: find better description
func (workerService Worker) Delete(worker atc.Worker) error {
	return nil
}
