package common

type ModelValidator interface {}

type FormValidator interface {
	validateForm(map[string]interface{}) error
}