package common

// Entity interfaces describes the base
type entity interface {
	SaveEntity(data interface{}) (interface{}, error)
	GetEntityById(id int) (entity, error)
}