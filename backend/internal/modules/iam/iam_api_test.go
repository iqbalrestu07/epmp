package iam_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/epmp/backend/internal/testutil"
)

type iamEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

func iamDo(t *testing.T, env *testutil.Env, method, path, body, token string) *http.Response {
	t.Helper()
	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}
	req, err := http.NewRequest(method, env.Server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func TestIAMAPI_RegisterLoginMe(t *testing.T) {
	env := testutil.Setup(t)
	email := fmt.Sprintf("iam_%d@epmp-test.com", time.Now().UnixNano())

	// Register
	resp := iamDo(t, env, "POST", "/api/v1/auth/register",
		fmt.Sprintf(`{"email":"%s","password":"secret123","name":"IAM Test"}`, email), "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d", resp.StatusCode)
	}

	// Login with wrong password → 401
	resp = iamDo(t, env, "POST", "/api/v1/auth/login",
		fmt.Sprintf(`{"email":"%s","password":"wrongpass"}`, email), "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong password login: expected 401, got %d", resp.StatusCode)
	}

	// Login correctly
	resp = iamDo(t, env, "POST", "/api/v1/auth/login",
		fmt.Sprintf(`{"email":"%s","password":"secret123"}`, email), "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login: expected 200, got %d", resp.StatusCode)
	}

	var envJson iamEnvelope
	var tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envJson); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	if err := json.Unmarshal(envJson.Data, &tokens); err != nil {
		t.Fatalf("unmarshal tokens: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("login returned empty tokens")
	}

	// /auth/me with the access token
	resp2 := iamDo(t, env, "GET", "/api/v1/auth/me", "", tokens.AccessToken)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("me: expected 200, got %d", resp2.StatusCode)
	}
	var meEnv iamEnvelope
	var me struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&meEnv); err != nil {
		t.Fatalf("decode me: %v", err)
	}
	if err := json.Unmarshal(meEnv.Data, &me); err != nil {
		t.Fatalf("unmarshal me: %v", err)
	}
	if me.Email != email {
		t.Errorf("me returned email %q, want %q", me.Email, email)
	}

	// Refresh rotates the access token
	resp3 := iamDo(t, env, "POST", "/api/v1/auth/refresh",
		fmt.Sprintf(`{"refresh_token":"%s"}`, tokens.RefreshToken), "")
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("refresh: expected 200, got %d", resp3.StatusCode)
	}

	// Logout invalidates the refresh token
	resp4 := iamDo(t, env, "POST", "/api/v1/auth/logout",
		fmt.Sprintf(`{"refresh_token":"%s"}`, tokens.RefreshToken), "")
	resp4.Body.Close()
	if resp4.StatusCode != http.StatusNoContent && resp4.StatusCode != http.StatusOK {
		t.Fatalf("logout: expected 204/200, got %d", resp4.StatusCode)
	}
	resp5 := iamDo(t, env, "POST", "/api/v1/auth/refresh",
		fmt.Sprintf(`{"refresh_token":"%s"}`, tokens.RefreshToken), "")
	resp5.Body.Close()
	if resp5.StatusCode != http.StatusUnauthorized {
		t.Fatalf("refresh after logout: expected 401, got %d", resp5.StatusCode)
	}
}

func TestIAMAPI_RegisterDuplicate(t *testing.T) {
	env := testutil.Setup(t)
	email := fmt.Sprintf("dup_%d@epmp-test.com", time.Now().UnixNano())
	body := fmt.Sprintf(`{"email":"%s","password":"secret123","name":"Dup"}`, email)

	resp := iamDo(t, env, "POST", "/api/v1/auth/register", body, "")
	resp.Body.Close()
	resp = iamDo(t, env, "POST", "/api/v1/auth/register", body, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate register: expected 409, got %d", resp.StatusCode)
	}
}

func TestIAMAPI_UsersRequirePermission(t *testing.T) {
	env := testutil.Setup(t)

	// Fresh user has no roles → /users must be forbidden.
	resp := env.Do(t, "GET", "/api/v1/users", "", env.Token, env.OrgID)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("users without permission: expected 403, got %d", resp.StatusCode)
	}
}

func TestIAMAPI_Unauthorized(t *testing.T) {
	env := testutil.Setup(t)
	testutil.AssertUnauthorized(t, env, "/api/v1/auth/me")
}
