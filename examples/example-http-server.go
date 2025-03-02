package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"

    rbac "github.com/asqewlabs/rbacteria"
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
    mux.Handle("GET /admin/dashboard", rbacManager.Middleware("view:admin.dashboard", adminDashboard))

    log.Fatal(http.ListenAndServe("localhost:80", mux))

}
