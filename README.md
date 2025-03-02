# rbacteria

## A lightweight and flexible Role-Based Access Control (RBAC) library for Go, designed to handle hierarchical roles, permissions, and inheritance.

## Features

- ✅ Define roles and permissions in JSON
- ✅ Support for role inheritance
- ✅ Fast permission checks
- ✅ Simple and easy-to-use API
- ✅ No external dependencies

## Installation

```sh
go get github.com/asqewlabs/rbacteria
```

## Usage

### 1. Define Your Roles in JSON

Create a `roles.json` file with role definitions:

```json
{
    "hr_manager": {
        "permissions": ["view:admin.dashboard", "read:admin.settings", "write:admin.settings"],
        "inherits": []
    },
    "it_manager": {
        "permissions": ["view:it.dashboard", "read:system.settings"],
        "inherits": []
    },
    "sysadmin": {
        "inherits": ["hr_manager", "it_manager"],
        "permissions": ["write:system.settings", "view:sysadmin.dashboard"]
    }
}
```

### 2. Using rbacteria in a web application

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/asqewlabs/rbacteria"
)

func adminDashboard( w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Admin Dashboard")
}

func main() {
    // Initialize RBAC
    rbacManager := rbac.NewRBAC()
    
    //Load roles file
    if err := rbacManager.LoadJSONFile("/path/to/roles.json"); err != nil {
        panic(err)
    }
    
    //Optionally, define how you will extract roles from a request
    //(If you want to override the default, which you probably do)
    rbacManager.WithExtractor(func(req *http.Request)[]string{
        return strings.Split(req.Header.Get("Roles"), ",")
    })
    
    //Set up your server mux
    mux := http.NewServeMux()
    
    /*Working on example*
    
}
```

## API Reference

### `rbacteria.NewRBAC()`

Creates a new RBAC instance with default values

#### Example:
```go
    rbac := rbacteria.NewRBAC()
```

### `rbacteria.LoadJSONFile(filename string) error`

Loads role definitions from JSON file.

#### Example:
```go
    rbac := rbacteria.NewRBAC()
    rbac.LoadJSONFile("rbac.json")
```

### `rbacteria.WithExtractor(extractor func(req *http.Request) []string)`

Sets the function to use to identify user roles from a request

#### Example:
```go
    rbac = rbacteria.NewRBAC().WithExtractor(func(req *http.Request)[]string {
        return []string{"role1", "role2"}
    })
```

### `rbacteria.WithLogger(logger *log.Logger)`

Sets the logger to use to log errors, access attempts, etc

#### Example:
```go
    rbac := rbacteria.NewRBAC().WithLogger(&log.Logger{})
```

### `rbacteria.HasPermission(assignedRoles []string, requiredPermission string, visited map[string]bool) bool`

Checks if a user has the required permission to access a resource

#### Example:
```go
    rbac := rbacteria.NewRBAC()
    if rbac.HasPermission([]string{"roleA"}, "actionB:resourceB", make(map[string]bool)) {
        log.Println("Permission Granted")
    } else {
        log.Println("Permission Denied")
    }
```

### `rbacteria.Middleware(permission string) func(http.HandlerFunc) http.HandlerFunc`

Middleware to validate permissions at endpoints

#### Example:
```go
    //TODO
```

## License

This project is licensed under the [AGPL-3.0 license](https://www.gnu.org/licenses/agpl-3.0.en.html). 


## Contributing

Pull requests are welcome! Please open an issue to discuss any changes before submitting a PR.

## Author

Will

