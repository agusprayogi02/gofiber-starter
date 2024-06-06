package config

import (
	"github.com/casbin/casbin/v2"
	"log"
)

var (
	C = "create"
	R = "read"
	L = "list"
	U = "update"
	D = "delete"
)

func LoadPermissions() {
	e, err := casbin.NewEnforcer("./asset/rbac/model.conf", "./asset/rbac/policy.csv")
	if err != nil {
		return
	}
	// Load the policy from DB.
	if err = e.LoadPolicy(); err != nil {
		log.Println("LoadPolicy failed, err: ", err)
	}

	// Check the permission.
	has, err := e.Enforce("admin", "files", R)
	if err != nil {
		log.Println("Enforce failed, err: ", err)
	}
	if !has {
		_, err := e.AddPolicies([][]string{
			{"admin", "files", C},
			{"admin", "files", R},
			{"admin", "files", U},
			{"admin", "files", D},
			{"admin", "files", L},
		})
		if err != nil {
			return
		}
		log.Println("do not have permission")
	}

	// Save the policy back to DB.
	if err = e.SavePolicy(); err != nil {
		log.Println("SavePolicy failed, err: ", err)
	}
}
