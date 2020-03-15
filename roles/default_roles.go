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

	db := config.GetDB()
	for _, s := range d.Roles {
		// Convert struct to models
		role := s.convertToRole()

		// Save records to database
		db.Create(&role)
	}
}

func (r DefaultRole) convertToRole() Role {
	policies := make([]Policy, len(r.Policies))

	for i, policy := range r.Policies {
		go func (i int, p PolicyConfig) {
			policies[i] = Policy{
				Resource: p.Resource,
				Index:    p.Index,
				Create:   p.Create,
				Read:     p.Read,
				Update:   p.Update,
				Delete:   p.Delete,
			}
		} (i, policy)
	}

	return Role{
		Name:		r.Name,
		Authority:	r.Authority,
		Policies:	policies,
	}
}