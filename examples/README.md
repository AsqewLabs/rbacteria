# rbacteria - example http server

##### Run with 

```bash
  rbacteria/examples$ go run example-http-server.go
```

##### Test direct role access - should display "Admin Dashboard"

```bash
  rbacteria/examples$ curl -H "Roles: hr_manager" localhost:8087/admin/dashboard
```

##### Test access denied - Should display "Forbidden"
```bash
  rbacteria/examples$ curl -H "Roles: it_manager" localhost:8087/admin/dashboard
```

##### Test inherited role access - should display "Admin Dashboard"
```bash
 rbacteria/examples$ curl -H "Roles: sysadmin" localhost:8087/admin/dashboard
```
