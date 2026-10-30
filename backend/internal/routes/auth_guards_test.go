// Route-level access matrix: every registered endpoint is listed with who may
// call it, and each one is checked against the real AuthMiddleware/VerifyRole.
package routes

import (
	"fmt"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/sqsclient"
	"inspirate-consulting/internal/supabase"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type access int

const (
	public  access = iota // no session needed
	anyRole               // any logged-in user, no role guard
	guarded               // logged-in user with one of the listed roles
)

type routeAccess struct {
	method string
	path   string // OpenAPI template, e.g. /user/{id}
	access access
	roles  []string // for guarded routes
}

const (
	roleStudent   = "student"
	roleCounselor = "counselor"
)

// Every endpoint the app registers. TestRouteTable_CoversEveryRegisteredRoute
// fails when a route is added without an entry here, so new endpoints get
// their access rules tested by default.
var routeTable = []routeAccess{
	// public
	{http.MethodGet, "/api/v1/health", public, nil},
	{http.MethodPost, "/user/login", public, nil},

	// auth & users
	{http.MethodPost, "/user/logout", anyRole, nil},
	{http.MethodPatch, "/user/reset-password", anyRole, nil},
	{http.MethodPost, "/user/{supabase_id}/force-password-reset", guarded, []string{roleCounselor}},
	{http.MethodPost, "/user", guarded, []string{roleCounselor}},
	{http.MethodGet, "/user/{id}", guarded, []string{roleStudent, roleCounselor}},

	// greeting
	{http.MethodPost, "/greeting", guarded, []string{roleCounselor}},

	// essays & reviews
	{http.MethodGet, "/students/{student_id}/essays", guarded, []string{roleStudent, roleCounselor}},
	{http.MethodPost, "/essays", guarded, []string{roleStudent}},
	{http.MethodPatch, "/essays/{essay_id}/status", guarded, []string{roleStudent, roleCounselor}},
	{http.MethodGet, "/essay-groups/{essay_group_id}/essays", anyRole, nil},
	{http.MethodPatch, "/essays/{essay_id}/group", anyRole, nil},
	{http.MethodPost, "/essay-groups", anyRole, nil},
	{http.MethodGet, "/students/{student_id}/essay-groups", anyRole, nil},
	{http.MethodPost, "/reviews", guarded, []string{roleStudent}},
	{http.MethodPost, "/reviews/{id}/refund", guarded, []string{roleCounselor}},
	{http.MethodPost, "/reviews/{id}/complete", guarded, []string{roleCounselor}},
	{http.MethodGet, "/essays/{essay_id}/review-status", guarded, []string{roleStudent, roleCounselor}},

	// students
	{http.MethodPatch, "/students/{id}/review-balance", guarded, []string{roleCounselor}},
	{http.MethodGet, "/students/{id}", anyRole, nil},
	{http.MethodGet, "/students/fetchall", guarded, []string{roleCounselor}},
	{http.MethodPatch, "/students/assign", guarded, []string{roleCounselor}},

	// counselors
	{http.MethodGet, "/counselors/reviews", guarded, []string{roleCounselor}},
	{http.MethodGet, "/counselors/students", guarded, []string{roleCounselor}},
	{http.MethodGet, "/counselors/fetchall", guarded, []string{roleCounselor}},

	// colleges & applications
	{http.MethodPost, "/colleges", guarded, []string{roleCounselor}},
	{http.MethodGet, "/colleges/{id}", anyRole, nil},
	{http.MethodGet, "/colleges", anyRole, nil},
	{http.MethodPut, "/colleges/{id}", anyRole, nil},
	{http.MethodDelete, "/colleges/{id}", anyRole, nil},
	{http.MethodGet, "/college-search", anyRole, nil},
	{http.MethodPost, "/applications", guarded, []string{roleStudent}},
	{http.MethodGet, "/applications", guarded, []string{roleStudent}},
	{http.MethodDelete, "/applications/{id}", guarded, []string{roleStudent}},
	{http.MethodPut, "/applications/{id}", guarded, []string{roleStudent}},
	{http.MethodPut, "/applications/{id}/rank", guarded, []string{roleStudent}},

	// to-do items
	{http.MethodPost, "/todo-items", guarded, []string{roleStudent, roleCounselor}},
	{http.MethodGet, "/todo-items", guarded, []string{roleStudent, roleCounselor}},
	{http.MethodGet, "/todo-items/search", guarded, []string{roleStudent, roleCounselor}},
	{http.MethodPatch, "/todo-items/{id}", guarded, []string{roleStudent, roleCounselor}},

	// videos
	{http.MethodPost, "/videos", guarded, []string{roleCounselor}},
	{http.MethodGet, "/videos/lookup", guarded, []string{roleStudent, roleCounselor}},
	{http.MethodGet, "/videos", guarded, []string{roleStudent, roleCounselor}},

	// media
	{http.MethodPost, "/media", anyRole, nil},
	{http.MethodGet, "/media/{id}", anyRole, nil},
	{http.MethodGet, "/media", anyRole, nil},
	{http.MethodDelete, "/media/{id}", anyRole, nil},
	{http.MethodPost, "/media-access", anyRole, nil},
	{http.MethodDelete, "/media-access/{id}", anyRole, nil},
	{http.MethodGet, "/media-access", anyRole, nil},

	// extracurriculars
	{http.MethodPost, "/extracurriculars/{studentID}", anyRole, nil},
	{http.MethodGet, "/extracurriculars", anyRole, nil},
	{http.MethodPatch, "/extracurriculars/{id}", anyRole, nil},

	// chat
	{http.MethodPost, "/messages", anyRole, nil},
	{http.MethodGet, "/chats", anyRole, nil},
	{http.MethodGet, "/chats/{user_id}/messages", anyRole, nil},
	{http.MethodPost, "/chats/{user_id}/read", anyRole, nil},
	{http.MethodPatch, "/messages/{id}", anyRole, nil},
	{http.MethodPatch, "/messages/{id}/read", anyRole, nil},

	// notification preferences
	{http.MethodGet, "/notification-preferences", anyRole, nil},
	{http.MethodPut, "/notification-preferences", anyRole, nil},
	{http.MethodPost, "/queue/email-notification", anyRole, nil},
}

var pathParam = regexp.MustCompile(`\{[^}]+\}`)

// concrete fills path params; guards run before params are parsed, so any value works.
func (r routeAccess) concrete() string {
	return pathParam.ReplaceAllString(r.path, "00000000-0000-0000-0000-0000000000aa")
}

func (r routeAccess) String() string { return r.method + " " + r.path }

func (r routeAccess) allows(role string) bool {
	switch r.access {
	case public, anyRole:
		return true
	default:
		return slices.Contains(r.roles, role)
	}
}

// deniedByGuard reports whether a response is the auth layer turning the
// request away, as opposed to the handler running (whatever it then returns:
// 2xx, 404, 422 for an empty body, or an error from the empty test database).
// Handlers report their own 401s as {"code","message"}, so only huma's error
// model (which carries a detail) counts as the auth layer.
func deniedByGuard(resp apiResponse) bool {
	return (resp.Status == http.StatusUnauthorized && resp.Detail() != "") ||
		(resp.Status == http.StatusForbidden && resp.Detail() == "Insufficient role")
}

func TestRouteTable_CoversEveryRegisteredRoute(t *testing.T) {
	t.Parallel()

	signingKey()
	_, api, err := SetupApp(config.Config{Supabase: &supabase.MockSupabase{}}, &data.Repository{}, &sqsclient.SQSClient{})
	require.NoError(t, err)

	var registered []string
	for path, item := range api.OpenAPI().Paths {
		for method, op := range map[string]any{
			http.MethodGet: item.Get, http.MethodPost: item.Post, http.MethodPut: item.Put,
			http.MethodPatch: item.Patch, http.MethodDelete: item.Delete,
		} {
			if fmt.Sprint(op) != "<nil>" {
				registered = append(registered, method+" "+path)
			}
		}
	}

	var listed []string
	for _, r := range routeTable {
		listed = append(listed, r.String())
	}
	sort.Strings(registered)
	sort.Strings(listed)
	assert.Equal(t, registered, listed,
		"routeTable must list every registered route exactly once; add new endpoints with their expected access")
}

func TestGuards_ProtectedRoutesRequireLogin(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	for _, r := range routeTable {
		if r.access == public {
			continue
		}
		t.Run(r.String(), func(t *testing.T) {
			resp := a.do(r.method, r.concrete(), "", map[string]any{})
			assert.Equal(t, http.StatusUnauthorized, resp.Status, "body: %s", resp.Body)
			assert.Equal(t, "Token Not Found", resp.Detail())
		})
	}
}

func TestGuards_PublicRoutesNeedNoSession(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	for _, r := range routeTable {
		if r.access != public {
			continue
		}
		t.Run(r.String(), func(t *testing.T) {
			resp := a.do(r.method, r.concrete(), "", map[string]any{})
			assert.False(t, deniedByGuard(resp), "got %d: %s", resp.Status, resp.Body)
		})
	}
}

func TestGuards_ProtectedRoutesRejectLoggedOutSession(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	tok := tokenFor(roleCounselor)
	a.revoke(tok)

	for _, r := range routeTable {
		if r.access == public {
			continue
		}
		t.Run(r.String(), func(t *testing.T) {
			resp := a.do(r.method, r.concrete(), tok, map[string]any{})
			assert.Equal(t, http.StatusUnauthorized, resp.Status, "body: %s", resp.Body)
		})
	}
}

func TestGuards_RoleMatrix(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)

	for _, r := range routeTable {
		// a role the app doesn't know must get no more than authentication alone gives
		for _, role := range []string{roleStudent, roleCounselor, "admin"} {
			t.Run(r.String()+" as "+role, func(t *testing.T) {
				// fresh session per request: allowed calls like logout really end the session
				tok := tokenFor(role)
				resp := a.do(r.method, r.concrete(), tok, map[string]any{})
				if r.allows(role) {
					assert.False(t, deniedByGuard(resp),
						"%s should be allowed, got %d: %s", role, resp.Status, resp.Body)
				} else {
					assert.Equal(t, http.StatusForbidden, resp.Status, "body: %s", resp.Body)
					assert.Equal(t, "Insufficient role", resp.Detail())
				}
			})
		}
	}
}

// Each guarded route must actually exclude someone; a guard listing both roles
// still needs a valid role claim, which the matrix checks via "admin".
func TestGuards_CounselorOnlyAndStudentOnlyRoutes(t *testing.T) {
	t.Parallel()

	var counselorOnly, studentOnly []string
	for _, r := range routeTable {
		if r.access != guarded {
			continue
		}
		switch {
		case slices.Equal(r.roles, []string{roleCounselor}):
			counselorOnly = append(counselorOnly, r.String())
		case slices.Equal(r.roles, []string{roleStudent}):
			studentOnly = append(studentOnly, r.String())
		}
	}

	// pin the security-relevant sets so a guard change is a deliberate test change
	assert.ElementsMatch(t, []string{
		"POST /user",
		"POST /user/{supabase_id}/force-password-reset",
		"POST /greeting",
		"POST /reviews/{id}/refund",
		"POST /reviews/{id}/complete",
		"PATCH /students/{id}/review-balance",
		"POST /colleges",
		"POST /videos",
		"GET /counselors/reviews",
		"GET /counselors/students",
		"GET /counselors/fetchall",
		"GET /students/fetchall",
		"PATCH /students/assign",
	}, counselorOnly)
	assert.ElementsMatch(t, []string{
		"POST /essays",
		"POST /reviews",
		"POST /applications",
		"GET /applications",
		"DELETE /applications/{id}",
		"PUT /applications/{id}",
		"PUT /applications/{id}/rank",
	}, studentOnly)
}
