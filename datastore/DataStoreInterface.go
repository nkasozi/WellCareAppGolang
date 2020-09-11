package datastore

import (
	"gitlab.com/capslock-ltd/reconciler/backend-golang/datastore/Entities"
)


type DStoreInterface interface {

	Add(item Entities.EntityInterface) (string, error)

	Update(item Entities.EntityInterface) (string, error)

	GetById(item Entities.EntityInterface) (Entities.EntityInterface, error)

	GetAll(afterFetchFilter func(item Entities.EntityInterface) bool, dataStoreQuery string) ([]Entities.EntityInterface, error)
}
