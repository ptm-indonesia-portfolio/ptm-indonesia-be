//go:build integration

package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"ptm-indonesia/bootstrap"
	"ptm-indonesia/config"
	"ptm-indonesia/model"

	"github.com/stretchr/testify/require"
)

func TestUserCRUDRequiresLoginAndSupportsPagingIntegration(t *testing.T) {
	cfg, dbCleanup, dbTeardown := setupIntegrationEnvironment(t)
	defer dbCleanup()
	defer dbTeardown()

	issuerURL, closeOIDCServer := startMockGoogleOIDCServer(t, "integration-google-client", model.GoogleIdentity{
		Subject:       "google-user-crud",
		Email:         cfg.Admin.Email,
		EmailVerified: true,
		Name:          "Google Admin",
		Picture:       stringPointer("https://example.com/avatar.png"),
	})
	defer closeOIDCServer()

	t.Setenv("GOOGLE_CLIENT_ID", "integration-google-client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "integration-google-secret")
	t.Setenv("GOOGLE_ISSUER_URL", issuerURL)
	t.Setenv("GOOGLE_REDIRECT_URI", "http://localhost:9101/api/v1/auth/google/callback")
	t.Setenv("FRONTEND_URL", "http://localhost:3100/home")

	app, cleanup, err := bootstrap.InitializeHTTPApplication()
	require.NoError(t, err)
	defer cleanup()

	healthRequest := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	healthResponse, err := app.App.Test(healthRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, healthResponse.StatusCode)

	preflightRequest := httptest.NewRequest(http.MethodOptions, "/api/v1/users", nil)
	preflightRequest.Header.Set("Origin", "http://localhost:3100")
	preflightRequest.Header.Set("Access-Control-Request-Method", http.MethodPut)
	preflightRequest.Header.Set("Access-Control-Request-Headers", "Content-Type")

	preflightResponse, err := app.App.Test(preflightRequest)
	require.NoError(t, err)
	require.NotEqual(t, http.StatusNotFound, preflightResponse.StatusCode)
	require.Equal(t, "http://localhost:3100", preflightResponse.Header.Get("Access-Control-Allow-Origin"))
	require.Contains(t, preflightResponse.Header.Get("Access-Control-Allow-Methods"), http.MethodPut)
	require.Contains(t, preflightResponse.Header.Get("Access-Control-Allow-Methods"), http.MethodDelete)

	unauthorizedRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=10", nil)
	unauthorizedResponse, err := app.App.Test(unauthorizedRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, unauthorizedResponse.StatusCode)

	var unauthorizedPayload model.ErrorResponse
	require.NoError(t, json.NewDecoder(unauthorizedResponse.Body).Decode(&unauthorizedPayload))
	require.NotEmpty(t, unauthorizedPayload.Errors)

	authCookie := authenticateAdminAndGetAuthCookie(t, cfg, app)
	currentUser := requestCurrentUser(t, app, authCookie)

	invalidStatusBody := map[string]any{
		"name":    "Invalid Status User",
		"email":   "invalid.status@example.com",
		"address": "Jakarta",
		"telp":    "081200000000",
		"status":  9,
	}

	invalidStatusResponse := performJSONRequest(t, app, http.MethodPost, "/api/v1/users", invalidStatusBody, authCookie)
	require.Equal(t, http.StatusBadRequest, invalidStatusResponse.StatusCode)

	var invalidStatusPayload model.ErrorResponse
	require.NoError(t, json.NewDecoder(invalidStatusResponse.Body).Decode(&invalidStatusPayload))
	require.NotEmpty(t, invalidStatusPayload.Errors)

	invalidSortRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=10&sort_by=unknown&sort_direction=desc", nil)
	invalidSortRequest.AddCookie(authCookie)
	invalidSortResponse, err := app.App.Test(invalidSortRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, invalidSortResponse.StatusCode)

	var invalidSortPayload model.ErrorResponse
	require.NoError(t, json.NewDecoder(invalidSortResponse.Body).Decode(&invalidSortPayload))
	require.NotEmpty(t, invalidSortPayload.Errors)

	createBody := map[string]any{
		"name":    "Member Baru",
		"email":   "MEMBER.BARU@example.com",
		"address": "Jakarta Selatan",
		"telp":    "081234567890",
		"status":  model.UserStatusFreeMember,
	}

	createResponse := performJSONRequest(t, app, http.MethodPost, "/api/v1/users", createBody, authCookie)
	require.Equal(t, http.StatusCreated, createResponse.StatusCode)

	var createPayload struct {
		Message string             `json:"message"`
		Data    model.UserResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(createResponse.Body).Decode(&createPayload))
	require.Equal(t, "member.baru@example.com", createPayload.Data.Email)
	require.Equal(t, currentUser.Data.ID, createPayload.Data.CreatedBy)
	require.Equal(t, currentUser.Data.ID, createPayload.Data.UpdatedBy)
	require.Equal(t, model.UserStatusFreeMember, createPayload.Data.Status)

	duplicateCreateBody := map[string]any{
		"name":    "Member Duplicate",
		"email":   "member.baru@example.com",
		"address": "Jakarta Selatan",
		"telp":    "081200000001",
		"status":  model.UserStatusFreeMember,
	}

	duplicateCreateResponse := performJSONRequest(t, app, http.MethodPost, "/api/v1/users", duplicateCreateBody, authCookie)
	require.Equal(t, http.StatusBadRequest, duplicateCreateResponse.StatusCode)

	var duplicateCreatePayload model.ErrorResponse
	require.NoError(t, json.NewDecoder(duplicateCreateResponse.Body).Decode(&duplicateCreatePayload))
	require.NotEmpty(t, duplicateCreatePayload.Errors)

	secondCreateBody := map[string]any{
		"name":    "Member Kedua",
		"email":   "member.kedua@example.com",
		"address": "Surabaya",
		"telp":    "082200000000",
		"status":  model.UserStatusFreeMember,
	}

	secondCreateResponse := performJSONRequest(t, app, http.MethodPost, "/api/v1/users", secondCreateBody, authCookie)
	require.Equal(t, http.StatusCreated, secondCreateResponse.StatusCode)

	var secondCreatePayload struct {
		Message string             `json:"message"`
		Data    model.UserResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(secondCreateResponse.Body).Decode(&secondCreatePayload))
	require.Equal(t, "member.kedua@example.com", secondCreatePayload.Data.Email)
	require.Equal(t, model.UserStatusFreeMember, secondCreatePayload.Data.Status)

	listRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=10", nil)
	listRequest.AddCookie(authCookie)
	listResponse, err := app.App.Test(listRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, listResponse.StatusCode)

	var listPayload struct {
		Message string                 `json:"message"`
		Data    model.UserListResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(listResponse.Body).Decode(&listPayload))
	require.Equal(t, int64(2), listPayload.Data.Meta.TotalItems)
	require.Equal(t, 1, listPayload.Data.Meta.Page)
	require.Equal(t, 10, listPayload.Data.Meta.Limit)
	require.Equal(t, 1, listPayload.Data.Meta.TotalPages)
	require.Len(t, listPayload.Data.Items, 2)
	require.Equal(t, secondCreatePayload.Data.ID, listPayload.Data.Items[0].ID)
	require.Equal(t, createPayload.Data.ID, listPayload.Data.Items[1].ID)
	for _, item := range listPayload.Data.Items {
		require.NotEqual(t, cfg.Admin.Email, item.Email)
	}

	searchRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=10&search=member.baru", nil)
	searchRequest.AddCookie(authCookie)
	searchResponse, err := app.App.Test(searchRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, searchResponse.StatusCode)

	var searchPayload struct {
		Message string                 `json:"message"`
		Data    model.UserListResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(searchResponse.Body).Decode(&searchPayload))
	require.Equal(t, int64(1), searchPayload.Data.Meta.TotalItems)
	require.Len(t, searchPayload.Data.Items, 1)
	require.Equal(t, createPayload.Data.ID, searchPayload.Data.Items[0].ID)

	sortRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=10&sort_by=status&sort_direction=desc", nil)
	sortRequest.AddCookie(authCookie)
	sortResponse, err := app.App.Test(sortRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, sortResponse.StatusCode)

	var sortPayload struct {
		Message string                 `json:"message"`
		Data    model.UserListResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(sortResponse.Body).Decode(&sortPayload))
	require.Equal(t, int64(2), sortPayload.Data.Meta.TotalItems)
	require.Len(t, sortPayload.Data.Items, 2)
	require.Equal(t, secondCreatePayload.Data.ID, sortPayload.Data.Items[0].ID)
	require.Equal(t, createPayload.Data.ID, sortPayload.Data.Items[1].ID)
	for _, item := range sortPayload.Data.Items {
		require.NotEqual(t, cfg.Admin.Email, item.Email)
	}

	detailRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+uintToString(createPayload.Data.ID), nil)
	detailRequest.AddCookie(authCookie)
	detailResponse, err := app.App.Test(detailRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, detailResponse.StatusCode)

	var detailPayload struct {
		Message string             `json:"message"`
		Data    model.UserResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(detailResponse.Body).Decode(&detailPayload))
	require.Equal(t, createPayload.Data.ID, detailPayload.Data.ID)
	require.Equal(t, "Member Baru", detailPayload.Data.Name)

	updateBody := map[string]any{
		"name":    "Member Premium",
		"email":   "premium.member@example.com",
		"address": "Bandung",
		"telp":    "081111111111",
		"status":  model.UserStatusPremiumMember,
	}

	updateResponse := performJSONRequest(t, app, http.MethodPut, "/api/v1/users/"+uintToString(createPayload.Data.ID), updateBody, authCookie)
	require.Equal(t, http.StatusOK, updateResponse.StatusCode)

	var updatePayload struct {
		Message string             `json:"message"`
		Data    model.UserResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(updateResponse.Body).Decode(&updatePayload))
	require.Equal(t, "Member Premium", updatePayload.Data.Name)
	require.Equal(t, "premium.member@example.com", updatePayload.Data.Email)
	require.Equal(t, model.UserStatusPremiumMember, updatePayload.Data.Status)
	require.Equal(t, currentUser.Data.ID, updatePayload.Data.UpdatedBy)

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+uintToString(createPayload.Data.ID), nil)
	deleteRequest.AddCookie(authCookie)
	deleteResponse, err := app.App.Test(deleteRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, deleteResponse.StatusCode)

	var deletePayload struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(deleteResponse.Body).Decode(&deletePayload))
	require.NotEmpty(t, deletePayload.Message)

	gormDB, dbCleanup, err := config.NewDatabase(cfg)
	require.NoError(t, err)
	defer dbCleanup()

	var deletedUser model.User
	require.NoError(t, gormDB.Where("id = ?", createPayload.Data.ID).First(&deletedUser).Error)
	require.NotNil(t, deletedUser.DeletedAt)
	require.Nil(t, deletedUser.StatusRow)

	deletedDetailRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+uintToString(createPayload.Data.ID), nil)
	deletedDetailRequest.AddCookie(authCookie)
	deletedDetailResponse, err := app.App.Test(deletedDetailRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, deletedDetailResponse.StatusCode)

	var deletedDetailPayload model.ErrorResponse
	require.NoError(t, json.NewDecoder(deletedDetailResponse.Body).Decode(&deletedDetailPayload))
	require.NotEmpty(t, deletedDetailPayload.Errors)

	recreateBody := map[string]any{
		"name":    "Member Recreated",
		"email":   updatePayload.Data.Email,
		"address": "Yogyakarta",
		"telp":    "083300000000",
		"status":  model.UserStatusFreeMember,
	}

	recreateResponse := performJSONRequest(t, app, http.MethodPost, "/api/v1/users", recreateBody, authCookie)
	require.Equal(t, http.StatusCreated, recreateResponse.StatusCode)

	var recreatePayload struct {
		Message string             `json:"message"`
		Data    model.UserResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(recreateResponse.Body).Decode(&recreatePayload))
	require.Equal(t, updatePayload.Data.Email, recreatePayload.Data.Email)
	require.NotEqual(t, createPayload.Data.ID, recreatePayload.Data.ID)
	require.Equal(t, model.UserStatusFreeMember, recreatePayload.Data.Status)

	finalListRequest := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=10", nil)
	finalListRequest.AddCookie(authCookie)
	finalListResponse, err := app.App.Test(finalListRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, finalListResponse.StatusCode)

	var finalListPayload struct {
		Message string                 `json:"message"`
		Data    model.UserListResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(finalListResponse.Body).Decode(&finalListPayload))
	require.Equal(t, int64(2), finalListPayload.Data.Meta.TotalItems)
	require.Len(t, finalListPayload.Data.Items, 2)
	require.Equal(t, recreatePayload.Data.Email, finalListPayload.Data.Items[0].Email)
	require.Equal(t, secondCreatePayload.Data.Email, finalListPayload.Data.Items[1].Email)
}

func TestUserServiceRequiresSuperAdminIntegration(t *testing.T) {
	cfg, dbCleanup, dbTeardown := setupIntegrationEnvironment(t)
	defer dbCleanup()
	defer dbTeardown()

	seedIntegrationUser(t, cfg, &model.User{
		Name:      "Free Member Access",
		Email:     "free.member.access@example.com",
		Status:    model.UserStatusFreeMember,
		CreatedBy: 0,
		UpdatedBy: 0,
	})

	issuerURL, closeOIDCServer := startMockGoogleOIDCServer(t, "integration-google-client", model.GoogleIdentity{
		Subject:       "google-user-member-access",
		Email:         "free.member.access@example.com",
		EmailVerified: true,
		Name:          "Free Member Access",
	})
	defer closeOIDCServer()

	t.Setenv("GOOGLE_CLIENT_ID", "integration-google-client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "integration-google-secret")
	t.Setenv("GOOGLE_ISSUER_URL", issuerURL)
	t.Setenv("GOOGLE_REDIRECT_URI", "http://localhost:9101/api/v1/auth/google/callback")
	t.Setenv("FRONTEND_URL", "http://localhost:3100/home")

	app, cleanup, err := bootstrap.InitializeHTTPApplication()
	require.NoError(t, err)
	defer cleanup()

	authCookie := authenticateUserAndGetAuthCookie(t, cfg.Auth.CookieName, app, "free.member.access@example.com")

	request := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&limit=10", nil)
	request.AddCookie(authCookie)

	response, err := app.App.Test(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, response.StatusCode)

	var payload model.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	require.NotEmpty(t, payload.Errors)
}

func authenticateAdminAndGetAuthCookie(t *testing.T, cfg *config.AppConfig, app *bootstrap.HTTPApplication) *http.Cookie {
	t.Helper()

	return authenticateUserAndGetAuthCookie(t, cfg.Auth.CookieName, app, cfg.Admin.Email)
}

func authenticateUserAndGetAuthCookie(t *testing.T, cookieName string, app *bootstrap.HTTPApplication, expectedEmail string) *http.Cookie {
	t.Helper()

	stateCookie := prepareGoogleLoginState(t, app)

	callbackRequest := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/auth/google/callback?code=test-auth-code&state="+url.QueryEscape(stateCookie.Value),
		nil,
	)
	callbackRequest.AddCookie(stateCookie)

	callbackResponse, err := app.App.Test(callbackRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, callbackResponse.StatusCode)
	require.Contains(t, callbackResponse.Header.Get("Location"), "auth_status=success")

	authCookie := findCookieByName(callbackResponse.Cookies(), cookieName)
	require.NotNil(t, authCookie)

	mePayload := requestCurrentUser(t, app, authCookie)
	require.Equal(t, expectedEmail, mePayload.Data.Email)

	return authCookie
}

func performJSONRequest(t *testing.T, app *bootstrap.HTTPApplication, method string, path string, body any, authCookie *http.Cookie) *http.Response {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	if authCookie != nil {
		request.AddCookie(authCookie)
	}

	response, err := app.App.Test(request)
	require.NoError(t, err)

	return response
}

func uintToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}
