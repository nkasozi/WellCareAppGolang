package Entities

type EntityInterface interface {
	GetEntityId() string
	SetEntityId(Id string) EntityInterface
	GetDateCreated() string
	GetDateModified() string
}
