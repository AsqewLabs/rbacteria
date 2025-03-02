# rbacteria

##A lightweight and flexible Role-Based Access Control (RBAC) library for Go, designed to handle hierarchical roles, permissions, and inheritance.

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

### 2. Load the Roles in Your Application

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
    
    //Set up your routes, with RBAC Middleware
    mux.Handle("GET /admin/dashboard", rbacManager.Middleware(func (w http.Writer, req *http.Request) {
        
    }))
    
    log.Fatal(http.ListenAndServe("localhost:80", mux))

    // Example permission check
    role := "sysadmin"
    permission := "write:system.settings"

    if rbacManager.Can(role, permission) {
        fmt.Printf("%s has permission to %s\n", role, permission)
    } else {
        fmt.Printf("%s is NOT allowed to %s\n", role, permission)
    }
}
```

### 3. Using RBAC in a Web API (Example with Gorilla Mux)

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/yourusername/rbac"
)

var rbacManager *rbac.RBAC

func checkPermissionHandler(w http.ResponseWriter, r *http.Request) {
    role := r.URL.Query().Get("role")
    action := r.URL.Query().Get("action")

    response := map[string]bool{"allowed": rbacManager.Can(role, action)}
    json.NewEncoder(w).Encode(response)
}

func main() {
    jsonData, err := os.ReadFile("roles.json")
    if err != nil {
        log.Fatalf("Failed to read roles.json: %v", err)
    }

    rbacManager = rbac.NewRBAC()
    if err := rbacManager.LoadJSON(jsonData); err != nil {
        log.Fatalf("Failed to load RBAC roles: %v", err)
    }

    router := mux.NewRouter()
    router.HandleFunc("/check-permission", checkPermissionHandler).Methods("GET")

    log.Println("RBAC API running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

#### Test the API

```sh
curl "http://localhost:8080/check-permission?role=sysadmin&action=write:system.settings"
```

## API Reference

### `rbac.NewRBAC()`

Creates a new RBAC instance.

### `rbac.LoadJSON(data []byte) error`

Loads role definitions from JSON data.

### `rbac.Can(role, permission string) bool`

Checks if a role has a specific permission.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Contributing

Pull requests are welcome! Please open an issue to discuss any changes before submitting a PR.

## Author

[Your Name](https://github.com/yourusername)

Will

