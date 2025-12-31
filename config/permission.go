package config

import (
	"starter-gofiber/variables"

	"github.com/casbin/casbin/v2"
)

var (
	CREATE_P = "create"
	READ_P   = "read"
	UPDATE_P = "update"
	DELETE_P = "delete"
	LIST_P   = "list"
	Enforcer *casbin.Enforcer
)

func InitializePermission(enforcer *casbin.Enforcer) error {
	Enforcer = enforcer
	post, files := "post", "files"
	enforcer.AddPermissionForUser(variables.ADMIN_ROLE, post, CREATE_P)
	enforcer.AddPermissionForUser(variables.ADMIN_ROLE, post, READ_P)
	enforcer.AddPermissionForUser(variables.ADMIN_ROLE, post, UPDATE_P)
	enforcer.AddPermissionForUser(variables.ADMIN_ROLE, post, DELETE_P)
	enforcer.AddPermissionForUser(variables.ADMIN_ROLE, post, LIST_P)

	enforcer.AddPermissionForUser(variables.ADMIN_ROLE, files, READ_P)
	enforcer.SavePolicy()
	return enforcer.LoadPolicy()
}
