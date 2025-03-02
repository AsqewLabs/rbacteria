package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/asqewlabs/rbacteria"
)

func adminDashboard( w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Admin Dashboard")
}

/* Function Chain from https://gowebexamples.com/advanced-middleware/ */
func Chain(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
    for _, m := range middlewares {
        f = m(f)
    }
    return f
}


func main() {
    // Initialize RBAC
    rbacManager := rbacteria.NewRBAC()

    //Load roles file
    if err := rbacManager.LoadJSONFile("roles.json"); err != nil {
        panic(err)
    }

    //Optionally, define how you will extract roles from a request
    //(If you want to override the default, which you probably do)
    rbacManager.WithExtractor(func(req *http.Request)[]string{
        return strings.Split(req.Header.Get("Roles"), ",")
    }).WithLogger(&log.Logger{})

    //Set up your server mux
    mux := http.NewServeMux()

    //Set up your routes, with RBAC Middleware
    mux.Handle("GET /admin/dashboard", Chain(adminDashboard, rbacManager.Middleware("view:admin.dashboard")))

    log.Fatal(http.ListenAndServe("localhost:8087", mux))

}
