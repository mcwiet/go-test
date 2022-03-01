package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/casbin/casbin/v2"
	"github.com/mcwiet/go-test/pkg/authorization"
	"github.com/newbmiao/dynacasbin"
)

func main() {
	// Initialize a DynamoDB adapter and use it in a Casbin enforcer:
	//use aws credentials default
	config := &aws.Config{
		Region: aws.String("us-east-1"), // your region
	} // Your AWS configuration
	table := "casbin-rules"
	a, err := dynacasbin.NewAdapter(config, table) // Your aws configuration and data source.
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer("rbac_model.conf", a)
	if err != nil {
		panic(err)
	}

	authorizer := authorization.NewCasbinAuthorizer(e)

	// Modify the policy. autoSave is support
	// e.AddPolicy("admin", "data1", "read")
	// authorizer.AddPermission("alice", "data1", "read")
	authorizer.AddRoleForUser("alice", authorization.Admin)

	// authorizer.RemoveRoleForUser("alice", authorization.Admin)
	// authorizer.RemovePermission("alice", "data1", "read")

	// Check the permission.
	result, err := authorizer.IsAuthorized("alice", "data1", "read")
	if err != nil {
		panic(err)
	}
	if result {
		log.Println("alice can read data1")
	} else {
		log.Println("alice can not read data1")
	}
}
