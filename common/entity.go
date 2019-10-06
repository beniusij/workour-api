package common

// Entity interfaces describes the base
type entity interface {
	SaveEntity(data interface{}) (int, error)
	GetEntityById(id int) (entity, error)
}