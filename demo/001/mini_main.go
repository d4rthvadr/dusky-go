package main

import (
	"errors"
	"sync"
)

var ErrTruckNotFound = errors.New("truck not found")

type FleetManager interface {
	AddTruck(id string, cargo int) error
	GetTruck(id string) (Truck, error)
	RemoveTruck(id string) error
	UpdateTruckCargo(id string, cargo int) error
}

type truckManger struct {
	trucks map[string]*Truck
	sync.RWMutex
}

type ContextKey string

type Truck struct {
	ID    string
	Cargo int
}

func (tm *truckManger) AddTruck(id string, cargo int) error {
	tm.Lock()
	defer tm.Unlock()

	tm.trucks[id] = &Truck{ID: id, Cargo: cargo}
	return nil
}

func (tm *truckManger) GetTruck(id string) (Truck, error) {
	tm.RLock()
	defer tm.RUnlock()
	truck, exists := tm.trucks[id]
	if !exists {
		return Truck{}, ErrTruckNotFound
	}
	return *truck, nil
}

func (tm *truckManger) RemoveTruck(id string) error {
	tm.Lock()
	defer tm.Unlock()
	// skip existence check for brevity
	delete(tm.trucks, id)
	return nil
}

func (tm *truckManger) UpdateTruckCargo(id string, cargo int) error {
	tm.Lock()
	defer tm.Unlock()
	truck, exits := tm.trucks[id]
	if !exits {
		return ErrTruckNotFound
	}
	truck.Cargo = cargo
	return nil
}

func NewTruckManager() truckManger {
	return truckManger{
		trucks: make(map[string]*Truck),
	}
}
