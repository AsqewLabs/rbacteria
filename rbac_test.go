package rbacteria_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	rbac "asqew.io/rbacteria"
)

func assertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("got %q, want nil", err)
	}
}

func assertThrowsError(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func Test_PermissionToString(t *testing.T) {
	var permission rbac.Permission
	permission.Action = "action"
	permission.Resource = "resource"

	got := permission.String()
	want := "action:resource"

	assertStrings(t, got, want)

}

func Test_PermissionLoad(t *testing.T) {
	var permission rbac.Permission
	validPermissionString := "action:resource"
	invalid_badFormat := "invalidblah"
	invalid_noAction := ":resource"
	invalid_noResource := "action:"

	assertNoError(t, permission.Load(validPermissionString))
	assertStrings(t, permission.Action, "action")
	assertStrings(t, permission.Resource, "resource")

	assertThrowsError(t, permission.Load(invalid_badFormat).Error(), "invalid format")
	assertThrowsError(t, permission.Load(invalid_noAction).Error(), "no action specified")
	assertThrowsError(t, permission.Load(invalid_noResource).Error(), "no resource specified")

}

func Test_HasPermission(t *testing.T) {
	r := rbac.NewRBAC()
	r.LoadJSONFile("valid_rbac.json")

	tests := []struct {
		name       string
		roles      []string
		permission string
		want       bool
	}{
		{"Has direct permission", []string{"hr_manager"}, "read:admin.settings", true},
		{"Does not have permission", []string{"hr_manager"}, "read:system.settings", false},
		{"Has direct permission, multiple groups", []string{"hr_manager", "it_manager"}, "read:admin.settings", true},
		{"Has inherited permission", []string{"sysadmin"}, "read:admin.settings", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if r.HasPermission(tc.roles, tc.permission, make(map[string]bool)) != tc.want {
				t.Errorf("Expected %t got %t", tc.want, !tc.want)
			}
		})
	}
}

func Test_RBACMiddleware(t *testing.T) {
	r := rbac.NewRBAC()
	if err := r.LoadJSONFile("valid_rbac.json"); err != nil {
		return
	}

	// Define a dummy handler to be wrapped
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	getRoles := func(req *http.Request) []string {
		return strings.Split(req.Header.Get("Roles"), ",")
	}
	// Middleware that checks for "admin" permission
	adminDashboard := r.Middleware("view:admin.dashboard", getRoles)
	itDashboard := r.Middleware("view:it.dashboard", getRoles)
	sysadminDashboard := r.Middleware("view:sysadmin.dashboard", getRoles)
	inherited := r.Middleware("view:it.dashboard", getRoles)

	// Create test cases
	tests := []struct {
		name           string
		roles          []string
		expectedStatus int
		expectedBody   string
		middleware     func(http.HandlerFunc) http.HandlerFunc
	}{
		{"Unauthorized User", []string{"user"}, http.StatusForbidden, "Forbidden\n", adminDashboard},
		{"Authorized User", []string{"it_manager"}, http.StatusOK, "OK", itDashboard},
		{"Authorized User", []string{"sysadmin"}, http.StatusOK, "OK", sysadminDashboard},
		{"Inherited Permissions", []string{"sysadmin"}, http.StatusOK, "OK", inherited},
		{"Multiple Roles", []string{"hr_manager", "it_manager"}, http.StatusOK, "OK", adminDashboard},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Roles", strings.Join(tc.roles, ","))
			w := httptest.NewRecorder()

			// Apply middleware to handler and serve request
			tc.middleware(handler).ServeHTTP(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// Check response body
			if w.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func Test_LoadFile(t *testing.T) {
	assertThrowsError(t, rbac.NewRBAC().LoadJSONFile("doesnotexist").Error(), "open doesnotexist: no such file or directory")
	assertThrowsError(t, rbac.NewRBAC().LoadJSONFile("invalid_json.jsonx").Error(), "invalid character 'i' looking for beginning of value")
	assertThrowsError(t, rbac.NewRBAC().LoadJSONFile("invalid_rbac.json").Error(), "error loading permission for role it_manager: invalid format")
}
