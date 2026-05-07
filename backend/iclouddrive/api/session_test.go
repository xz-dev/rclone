package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractHeadersMergesCookies(t *testing.T) {
	s := NewSession(DefaultEndpoints())
	s.Cookies = []*http.Cookie{{Name: "existing", Value: "old"}}

	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Add("Set-Cookie", (&http.Cookie{Name: "existing", Value: "new"}).String())
	resp.Header.Add("Set-Cookie", (&http.Cookie{Name: "fresh", Value: "value"}).String())
	resp.Header.Set("X-Apple-Session-Token", "session-token")

	s.extractHeaders(resp)

	require.Len(t, s.Cookies, 2)
	assert.Equal(t, "new", s.Cookies[0].Value)
	assert.Equal(t, "fresh", s.Cookies[1].Name)
	assert.Equal(t, "session-token", s.SessionToken)
}

func TestExtractHeadersDeletesEmptyCookies(t *testing.T) {
	s := NewSession(DefaultEndpoints())
	s.Cookies = []*http.Cookie{{Name: "X-APPLE-WEBAUTH-HSA-LOGIN", Value: "stale"}, {Name: "keep", Value: "value"}}

	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Add("Set-Cookie", (&http.Cookie{Name: "X-APPLE-WEBAUTH-HSA-LOGIN", Value: ""}).String())

	s.extractHeaders(resp)

	require.Len(t, s.Cookies, 1)
	assert.Equal(t, "keep", s.Cookies[0].Name)
	assert.Equal(t, "keep=value", s.GetCookieString())
}

func TestChinaMainlandEndpointsUseGlobalAuth(t *testing.T) {
	endpoints := ChinaMainlandEndpoints()

	assert.Equal(t, "https://www.icloud.com.cn", endpoints.Base)
	assert.Equal(t, "https://setup.icloud.com.cn/setup/ws/1", endpoints.Setup)
	assert.Equal(t, "https://idmsa.apple.com/appleauth/auth", endpoints.Auth)
}

func TestSRPAuthHeadersUseRegionalRedirectURI(t *testing.T) {
	s := NewSession(ChinaMainlandEndpoints())
	s.ClientID = "client-id"

	headers := s.getSRPAuthHeaders()

	assert.Equal(t, "https://idmsa.apple.com", headers["Origin"])
	assert.Equal(t, "https://idmsa.apple.com/", headers["Referer"])
	assert.Equal(t, "https://www.icloud.com.cn", headers["X-Apple-OAuth-Redirect-URI"])
}
