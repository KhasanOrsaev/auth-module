Auth module
==

Auth module is used for authenticate, authorize users. Also for creating
users, groups and roles.

Currently Auth module has support with basic/jwt/password authorization.

##Install

Get module with command:
```go
go get github.com/KhasanOrsaev/auth-module
```
 
##Quick start

Before start using require to run a migration. You can do it on init.
Example:
```go
package main
import (
    "github.com/KhasanOrsaev/auth-module"
    "github.com/KhasanOrsaev/auth-module/internal/migration"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var (
    db, _ = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
)

func init()  {  
    migration.Migrate(db)
}

func main()  {
 authClient,err := auth_module.NewAuthClient(db, auth_module.NullType)
}
``` 

##Usage

Auth model has models:
* User
* Group
* Role

to create a new instance of those you just need call:
```go
// create or replace user with current login
authClient.ApplyUser(login, password string, roles []*models.Role)

// create or replace rule with current name
authClient.ApplyRule(roleName string, scope string)

// create or replace group with current name
authClient.ApplyRule(name string, roles []*models.Role)
```
---
For authorize or authenticate you should create auth client with needed
auth type:
```go
authClient,err := auth_module.NewAuthClient(db, auth_module.BearerType)
```

* __BearerType__ - jwt
* __BasicType__ - basic auth
* __PasswordType__ - auth by password