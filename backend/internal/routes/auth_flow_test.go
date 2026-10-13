// End-to-end auth flows through the real router, middleware and handlers:
// login/logout, account creation by a counselor, first login with the temp
// password, password reset, and counselor-forced resets.
package routes

import (
	"errors"
	"inspirate-consulting/internal/models"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createdAccount struct {
	User         models.User `json:"user"`
	TempPassword string      `json:"temp_password"`
}

func (a *authApp) userPath(supabaseID uuid.UUID) string {
	a.t.Helper()
	u, ok := a.users.BySupabaseID(supabaseID)
	require.True(a.t, ok, "no users row for %s", supabaseID)
	return "/user/" + u.ID.String()
}

func TestFlow_LoginThenLogout(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	acct := a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
	profile := a.userPath(acct.ID)

	// login
	resp, tok := a.login("casey@example.com", "Password1!")
	require.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)
	require.NotEmpty(t, tok)

	var user models.User
	resp.JSON(t, &user)
	assert.Equal(t, "Casey Counselor", user.Name)
	assert.Equal(t, acct.ID, user.SupabaseID)
	assert.Nil(t, user.ResetTime)

	t.Run("login sets a hardened session cookie", func(t *testing.T) {
		c := resp.Cookie("jwt")
		require.NotNil(t, c)
		assert.True(t, c.HttpOnly, "HttpOnly keeps the token away from page JavaScript")
		assert.True(t, c.Secure)
		assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
		assert.Equal(t, "/", c.Path)
		assert.Equal(t, 3600, c.MaxAge, "cookie lifetime follows the Supabase session")
	})

	t.Run("login body never carries the tokens", func(t *testing.T) {
		body := string(resp.Body)
		assert.NotContains(t, body, tok)
		assert.NotContains(t, body, "access_token")
		assert.NotContains(t, body, "refresh_token")
	})

	// the cookie authenticates later requests
	require.Equal(t, http.StatusOK, a.do(http.MethodGet, profile, tok, nil).Status)

	// logout
	out := a.do(http.MethodPost, "/user/logout", tok, nil)
	require.Equal(t, http.StatusNoContent, out.Status, "body: %s", out.Body)

	t.Run("logout revokes the Supabase session it was called with", func(t *testing.T) {
		assert.Equal(t, []string{tok}, a.sb.LogoutCalls)
	})

	t.Run("logout expires the cookie with matching attributes", func(t *testing.T) {
		c := out.Cookie("jwt")
		require.NotNil(t, c)
		assert.Empty(t, c.Value)
		assert.Negative(t, c.MaxAge, "Max-Age=0 tells the browser to delete it")
		assert.Equal(t, "/", c.Path, "must match the login cookie's path to replace it")
		assert.True(t, c.HttpOnly)
		assert.True(t, c.Secure)
		assert.Contains(t, out.Header.Get("Set-Cookie"), "Max-Age=0")
	})

	t.Run("the old token is dead after logout", func(t *testing.T) {
		resp := a.do(http.MethodGet, profile, tok, nil)
		assert.Equal(t, http.StatusUnauthorized, resp.Status)
		assert.Equal(t, "Invalid/Expired Token", resp.Detail())
	})

	t.Run("logging out twice is rejected", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, a.do(http.MethodPost, "/user/logout", tok, nil).Status)
	})

	t.Run("logging back in issues a new working session", func(t *testing.T) {
		tok2 := a.mustLogin("casey@example.com", "Password1!")
		assert.NotEqual(t, tok, tok2)
		assert.Equal(t, http.StatusOK, a.do(http.MethodGet, profile, tok2, nil).Status)
	})
}

func TestFlow_LogoutRequiresSession(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	resp := a.do(http.MethodPost, "/user/logout", "", nil)
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
	assert.Nil(t, resp.Cookie("jwt"))
	assert.Empty(t, a.sb.LogoutCalls, "Supabase must not be called without a session")
}

func TestFlow_LogoutEndsOnlyThatSession(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	acct := a.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
	profile := a.userPath(acct.ID)

	laptop := a.mustLogin("sam@example.com", "Password1!")
	phone := a.mustLogin("sam@example.com", "Password1!")
	require.NotEqual(t, laptop, phone)

	require.Equal(t, http.StatusNoContent, a.do(http.MethodPost, "/user/logout", laptop, nil).Status)

	assert.Equal(t, http.StatusUnauthorized, a.do(http.MethodGet, profile, laptop, nil).Status)
	assert.Equal(t, http.StatusOK, a.do(http.MethodGet, profile, phone, nil).Status)
}

func TestFlow_SessionRevokedElsewhere(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	acct := a.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
	profile := a.userPath(acct.ID)
	tok := a.mustLogin("sam@example.com", "Password1!")

	// e.g. an admin signs the user out in Supabase; the JWT itself is still unexpired
	a.sb.Revoke(tok)

	assert.Equal(t, http.StatusUnauthorized, a.do(http.MethodGet, profile, tok, nil).Status)
}

func TestFlow_LoginFailures(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	a.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
	// an auth account whose users row is missing (e.g. a half-finished signup)
	a.sb.AddAccount("orphan@example.com", "Password1!", roleStudent)

	tests := []struct {
		name   string
		body   any
		status int
	}{
		{"wrong password", map[string]string{"email": "sam@example.com", "password": "nope"}, http.StatusBadRequest},
		{"password with different case", map[string]string{"email": "sam@example.com", "password": "password1!"}, http.StatusBadRequest},
		{"unknown email", map[string]string{"email": "ghost@example.com", "password": "Password1!"}, http.StatusBadRequest},
		{"no users row for the auth account", map[string]string{"email": "orphan@example.com", "password": "Password1!"}, http.StatusNotFound},
		{"missing password", map[string]string{"email": "sam@example.com"}, http.StatusUnprocessableEntity},
		{"missing email", map[string]string{"password": "Password1!"}, http.StatusUnprocessableEntity},
		{"empty body", map[string]string{}, http.StatusUnprocessableEntity},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := a.do(http.MethodPost, "/user/login", "", tc.body)
			assert.Equal(t, tc.status, resp.Status, "body: %s", resp.Body)
			assert.Nil(t, resp.Cookie("jwt"), "a failed login must not set a session cookie")
		})
	}

	t.Run("supabase unavailable", func(t *testing.T) {
		b := newAuthApp(t, nil)
		b.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
		b.sb.LoginErr = errors.New("dial tcp: connection refused")

		resp, tok := b.login("sam@example.com", "Password1!")
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Empty(t, tok)
	})
}

// A counselor creates a student; the student logs in with the temp password,
// is told to reset, resets, and from then on logs in with the new password and
// only reaches student-allowed endpoints.
func TestFlow_CounselorCreatesStudent_FirstLoginAndReset(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
	counselorTok := a.mustLogin("casey@example.com", "Password1!")

	// 1. counselor creates the account
	resp := a.do(http.MethodPost, "/user", counselorTok, map[string]any{
		"name": "New Student", "email": "new.student@example.com", "pfp_key": nil, "role": "student",
	})
	require.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)

	var created createdAccount
	resp.JSON(t, &created)
	assert.Equal(t, "New Student", created.User.Name)
	require.NotEmpty(t, created.TempPassword)
	assert.NotNil(t, created.User.ResetTime, "new accounts must start with a pending password reset")

	acct, ok := a.sb.Account("new.student@example.com")
	require.True(t, ok, "Supabase account was not created")
	assert.Equal(t, roleStudent, acct.Role, "role is stored in Supabase app_metadata")
	assert.Equal(t, created.TempPassword, acct.Password)
	assert.Equal(t, acct.ID, created.User.SupabaseID)

	// 2. first login with the temp password reports the pending reset
	resp, studentTok := a.login("new.student@example.com", created.TempPassword)
	require.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)
	var first models.User
	resp.JSON(t, &first)
	assert.NotNil(t, first.ResetTime, "login must surface reset_time so the frontend can redirect")

	// 3. the student resets their password
	reset := a.do(http.MethodPatch, "/user/reset-password", studentTok, map[string]string{"new_password": "MyOwnPass1!"})
	require.Equal(t, http.StatusNoContent, reset.Status, "body: %s", reset.Body)

	require.Len(t, a.sb.ResetPasswordCalls, 1)
	assert.Equal(t, acct.ID.String(), a.sb.ResetPasswordCalls[0].UserID,
		"reset must target the caller's own account (Supabase-ID from the session)")
	stored, _ := a.users.BySupabaseID(acct.ID)
	assert.Nil(t, stored.ResetTime, "reset clears the pending flag")

	// 4. old password stops working, new one works and no reset is pending
	old, oldTok := a.login("new.student@example.com", created.TempPassword)
	assert.Equal(t, http.StatusBadRequest, old.Status)
	assert.Empty(t, oldTok)

	resp, studentTok = a.login("new.student@example.com", "MyOwnPass1!")
	require.Equal(t, http.StatusOK, resp.Status)
	var again models.User
	resp.JSON(t, &again)
	assert.Nil(t, again.ResetTime)

	// 5. the new account has student access only
	t.Run("student can reach student routes", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, a.do(http.MethodGet, a.userPath(acct.ID), studentTok, nil).Status)
		assert.False(t, deniedByGuard(a.do(http.MethodGet, "/applications", studentTok, nil)))
	})
	t.Run("student cannot reach counselor routes", func(t *testing.T) {
		for _, r := range []struct{ method, path string }{
			{http.MethodPost, "/user"},
			{http.MethodPost, "/colleges"},
			{http.MethodPost, "/user/" + acct.ID.String() + "/force-password-reset"},
		} {
			resp := a.do(r.method, r.path, studentTok, map[string]any{})
			assert.Equal(t, http.StatusForbidden, resp.Status, "%s %s", r.method, r.path)
		}
	})
}

func TestFlow_CounselorCreatesCounselor(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
	tok := a.mustLogin("casey@example.com", "Password1!")

	resp := a.do(http.MethodPost, "/user", tok, map[string]any{
		"name": "Second Counselor", "email": "second@example.com", "pfp_key": nil, "role": "counselor",
	})
	require.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)
	var created createdAccount
	resp.JSON(t, &created)

	// the new counselor's role comes through login into the token
	newTok := a.mustLogin("second@example.com", created.TempPassword)
	resp = a.do(http.MethodPost, "/user", newTok, map[string]any{
		"name": "Their Student", "email": "their.student@example.com", "pfp_key": nil, "role": "student",
	})
	assert.Equal(t, http.StatusOK, resp.Status, "a counselor-created counselor must be able to create accounts: %s", resp.Body)
}

func TestFlow_CreateAccountRejected(t *testing.T) {
	t.Parallel()

	valid := map[string]any{"name": "X", "email": "x@example.com", "pfp_key": nil, "role": "student"}

	t.Run("without a session", func(t *testing.T) {
		a := newAuthApp(t, nil)
		resp := a.do(http.MethodPost, "/user", "", valid)
		assert.Equal(t, http.StatusUnauthorized, resp.Status)
		assert.Empty(t, a.sb.SignupCalls)
		assert.Zero(t, a.users.Count())
	})

	t.Run("as a student", func(t *testing.T) {
		a := newAuthApp(t, nil)
		a.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
		tok := a.mustLogin("sam@example.com", "Password1!")

		resp := a.do(http.MethodPost, "/user", tok, valid)
		assert.Equal(t, http.StatusForbidden, resp.Status)
		assert.Empty(t, a.sb.SignupCalls, "the guard must stop the request before Supabase signup")
		assert.Equal(t, 1, a.users.Count())
	})

	invalid := []struct {
		name string
		body map[string]any
	}{
		{"unknown role", map[string]any{"name": "X", "email": "x@example.com", "pfp_key": nil, "role": "admin"}},
		{"missing role", map[string]any{"name": "X", "email": "x@example.com", "pfp_key": nil}},
		{"empty name", map[string]any{"name": "", "email": "x@example.com", "pfp_key": nil, "role": "student"}},
		{"missing email", map[string]any{"name": "X", "pfp_key": nil, "role": "student"}},
	}
	for _, tc := range invalid {
		t.Run("invalid input: "+tc.name, func(t *testing.T) {
			a := newAuthApp(t, nil)
			a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
			tok := a.mustLogin("casey@example.com", "Password1!")

			resp := a.do(http.MethodPost, "/user", tok, tc.body)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.Status, "body: %s", resp.Body)
			assert.Empty(t, a.sb.SignupCalls)
			assert.Equal(t, 1, a.users.Count())
		})
	}

	t.Run("duplicate email", func(t *testing.T) {
		a := newAuthApp(t, nil)
		a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
		tok := a.mustLogin("casey@example.com", "Password1!")

		resp := a.do(http.MethodPost, "/user", tok, map[string]any{
			"name": "Dup", "email": "casey@example.com", "pfp_key": nil, "role": "student",
		})
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Status, "body: %s", resp.Body)
		assert.Equal(t, 1, a.users.Count(), "no users row for a failed signup")
		acct, _ := a.sb.Account("casey@example.com")
		assert.Equal(t, roleCounselor, acct.Role, "the existing account is untouched")
	})

	t.Run("supabase signup fails", func(t *testing.T) {
		a := newAuthApp(t, nil)
		a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
		tok := a.mustLogin("casey@example.com", "Password1!")
		a.sb.SignupErr = errors.New("supabase down")

		resp := a.do(http.MethodPost, "/user", tok, valid)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Equal(t, 1, a.users.Count())
		assert.False(t, strings.Contains(string(resp.Body), "temp_password"))
	})
}

func TestFlow_CounselorForcesPasswordReset(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
	stu := a.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
	counselorTok := a.mustLogin("casey@example.com", "Password1!")

	// no reset pending to begin with
	resp, _ := a.login("sam@example.com", "Password1!")
	var before models.User
	resp.JSON(t, &before)
	require.Nil(t, before.ResetTime)

	// counselor forces a reset
	force := a.do(http.MethodPost, "/user/"+stu.ID.String()+"/force-password-reset", counselorTok, nil)
	require.Equal(t, http.StatusNoContent, force.Status, "body: %s", force.Body)

	// the student's next login reports it
	resp, studentTok := a.login("sam@example.com", "Password1!")
	require.Equal(t, http.StatusOK, resp.Status)
	var after models.User
	resp.JSON(t, &after)
	require.NotNil(t, after.ResetTime, "forced reset must show up on the next login")

	// resetting clears it again
	reset := a.do(http.MethodPatch, "/user/reset-password", studentTok, map[string]string{"new_password": "Brand-New-1"})
	require.Equal(t, http.StatusNoContent, reset.Status, "body: %s", reset.Body)
	resp, _ = a.login("sam@example.com", "Brand-New-1")
	var cleared models.User
	resp.JSON(t, &cleared)
	assert.Nil(t, cleared.ResetTime)
}

func TestFlow_ForcePasswordResetRejected(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)
	victim := a.seedUser("Vic Student", "vic@example.com", "Password1!", roleStudent)
	a.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
	path := "/user/" + victim.ID.String() + "/force-password-reset"

	t.Run("without a session", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, a.do(http.MethodPost, path, "", nil).Status)
	})

	t.Run("by a student", func(t *testing.T) {
		tok := a.mustLogin("sam@example.com", "Password1!")
		assert.Equal(t, http.StatusForbidden, a.do(http.MethodPost, path, tok, nil).Status)
	})

	t.Run("with a malformed id", func(t *testing.T) {
		tok := a.mustLogin("casey@example.com", "Password1!")
		resp := a.do(http.MethodPost, "/user/not-a-uuid/force-password-reset", tok, nil)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Status)
	})

	stored, _ := a.users.BySupabaseID(victim.ID)
	assert.Nil(t, stored.ResetTime, "rejected requests must not flag the account")
}

func TestFlow_ResetPasswordRequiresSession(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	resp := a.do(http.MethodPatch, "/user/reset-password", "", map[string]string{"new_password": "Whatever1!"})
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
	assert.Empty(t, a.sb.ResetPasswordCalls)
}

// The role that counts is the one Supabase holds for the account, carried in
// the signed token; the same person can't pick a role at login.
func TestFlow_RoleFollowsTheAccount(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	a.seedUser("Sam Student", "sam@example.com", "Password1!", roleStudent)
	a.seedUser("Casey Counselor", "casey@example.com", "Password1!", roleCounselor)

	studentTok := a.mustLogin("sam@example.com", "Password1!")
	counselorTok := a.mustLogin("casey@example.com", "Password1!")

	for _, r := range routeTable {
		if r.access != guarded {
			continue
		}
		t.Run(r.String(), func(t *testing.T) {
			assert.Equal(t, r.allows(roleStudent), !deniedByGuard(a.do(r.method, r.concrete(), studentTok, map[string]any{})), "student")
			assert.Equal(t, r.allows(roleCounselor), !deniedByGuard(a.do(r.method, r.concrete(), counselorTok, map[string]any{})), "counselor")
		})
	}
}
