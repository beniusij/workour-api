package common

// Entity interfaces describes the base
type entity interface {
	Save(data interface{}) (int, error)
	GetById(id int) (entity, error)
}