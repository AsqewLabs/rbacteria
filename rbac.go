package rbacteria

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type RBAC struct {
	Roles  map[string]Role
	Logger *log.Logger
	RoleExtractor func(r *http.Request) []string
}

func NewRBAC() *RBAC {
	return &RBAC{
		Roles:  make(map[string]Role),
		Logger: log.Default(),
		RoleExtractor: func(req *http.Request) []string {
			return strings.Split(req.Header.Get("Roles"), ",")
		},
	}
}


func (r *RBAC) WithLogger(logger *log.Logger) *RBAC {
	r.Logger = logger
	return r
}

func (r *RBAC) WithExtractor(extractor func(req *http.Request) []string) *RBAC {
	r.RoleExtractor = extractor
	return r
}

func (r *RBAC) Middleware(permission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			roles := r.RoleExtractor(req)
			if r.HasPermission(roles, permission, make(map[string]bool)) {
				r.Logger.Printf("granting access to %s from permission %s via role %s", req.URL, permission, roles)
				f(w, req)
			} else {
				r.Logger.Printf("denied access to %s due to missing permission %s", req.URL, permission)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			}
		}
	}
}

func (r *RBAC) HasPermission(assignedRoles []string, requiredPermission string, visited map[string]bool) bool {
	permitted := false
	for _, assigned := range assignedRoles {
		if !visited[assigned] {
			if role, ok := r.Roles[assigned]; !ok {
				//role does not exist, so skip it
				continue
			} else {
				for _, permission := range role.Permissions {
					if permission.String() == requiredPermission {
						//Permission has been found! Set permitted to true and bust out of the loop
						permitted = true
						break
					}
				}
				if !permitted {
					//check inherited roles
					permitted = r.HasPermission(role.Inherits, requiredPermission, visited)
				}
			}
			if permitted {
				//Don't need to continue checking other roles
				break
			}
		}
	}

	return permitted
}

type roleJSON struct {
	Permissions []string `json:"permissions"`
	Inherits    []string `json:"inherits"`
}

func (r *RBAC) LoadJSONFile(filename string) error {
	if content, err := os.ReadFile(filename); err != nil {
		return err
	} else {
		var roleMap map[string]roleJSON
		if err := json.Unmarshal(content, &roleMap); err != nil {
			return err
		}

		for roleName, roleData := range roleMap {
			role := Role{
				Name:        roleName,
				Inherits:    roleData.Inherits,
				Permissions: make([]Permission, 0, len(roleData.Permissions)),
			}

			for _, permString := range roleData.Permissions {
				var perm Permission
				if err := perm.Load(permString); err != nil {
					return fmt.Errorf("error loading permission for role %s: %v", roleName, err)
				}
				role.Permissions = append(role.Permissions, perm)
			}

			r.Roles[roleName] = role
		}

		return nil
	}
}

type Role struct {
	Name        string       //Name of the role, to be assigned to a user
	Permissions []Permission //A role can have multiple permissions
	Inherits    []string     //Names of the roles it inherits from
}

type Permission struct {
	Action   string //Create, read, update, destroy, view
	Resource string //"dashboard", "invoice", etc
}

func (p *Permission) Load(input string) error {
	splitString := strings.SplitN(input, ":", 2)
	if len(splitString) < 2 {
		return errors.New(("invalid format"))
	}

	action := splitString[0]
	resource := splitString[1]

	if len(action) == 0 {
		return errors.New("no action specified")
	} else if len(resource) == 0 {
		return errors.New("no resource specified")
	}
	p.Action = action
	p.Resource = resource

	return nil
}

func (p *Permission) String() string {
	return p.Action + ":" + p.Resource
}
