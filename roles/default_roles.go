package roles

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"workour-api/config"
)

const filename = "roles/default_roles.yaml"

type PolicyConfig struct {
	Resource 	string `yaml:"Resource"`
	Index		bool `yaml:"Index"`
	Create		bool `yaml:"Create"`
	Read		bool `yaml:"Read"`
	Update		bool `yaml:"Update"`
	Delete		bool `yaml:"Delete"`
}

type DefaultRole struct {
	Name 		string `yaml:"Name"`
	Authority 	int `yaml:"Authority"`
	Policies 	[]PolicyConfig `yaml:"Policies"`
}

type DefaultRoles struct {
	Roles []DefaultRole `yaml:"Roles"`
}

func CreateDefaultRoles() {
	// Load default roles from config
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal into struct
	d := DefaultRoles{}
	err = yaml.Unmarshal(file, &d)
	if err != nil {
		log.Fatal(err)
	}

	// Save to database
	db := config.GetDB()

	for _, role := range d.Roles {
		db.Create(&role)
	}
}