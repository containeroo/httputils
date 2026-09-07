package httputils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	t.Parallel()

	t.Run("valid headers", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders([]string{
			"Content-Type=application/json",
			"Authorization=Bearer token",
		}, false)

		require.NoError(t, err)
		assert.Equal(t, http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer token"},
		}, result)
	})

	t.Run("canonicalizes header names", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders([]string{"content-type=application/json"}, false)

		require.NoError(t, err)
		assert.Equal(t, http.Header{"Content-Type": {"application/json"}}, result)
	})

	t.Run("empty headers", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders(nil, false)

		require.NoError(t, err)
		assert.Equal(t, http.Header{}, result)
	})

	t.Run("malformed header", func(t *testing.T) {
		t.Parallel()

		_, err := ParseHeaders([]string{"AuthorizationBearer token"}, false)

		assert.EqualError(t, err, "invalid header format: AuthorizationBearer token")
	})

	t.Run("header with spaces", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders([]string{"  Content-Type = application/json  "}, false)

		require.NoError(t, err)
		assert.Equal(t, http.Header{"Content-Type": {"application/json"}}, result)
	})

	t.Run("header with empty key", func(t *testing.T) {
		t.Parallel()

		_, err := ParseHeaders([]string{"=value"}, false)

		assert.EqualError(t, err, "header key cannot be empty: =value")
	})

	t.Run("header with empty value", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders([]string{"X-Empty="}, false)

		require.NoError(t, err)
		assert.Equal(t, http.Header{"X-Empty": {""}}, result)
	})

	t.Run("preserves comma in value", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders([]string{"Accept=text/html, application/json"}, false)

		require.NoError(t, err)
		assert.Equal(t, "text/html, application/json", result.Get("Accept"))
	})

	t.Run("preserves equals in value", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders([]string{"Authorization=Signature token=abc123"}, false)

		require.NoError(t, err)
		assert.Equal(t, "Signature token=abc123", result.Get("Authorization"))
	})

	t.Run("preserves duplicate headers when allowed", func(t *testing.T) {
		t.Parallel()

		result, err := ParseHeaders([]string{
			"X-Test=one",
			"X-Test=two",
		}, true)

		require.NoError(t, err)
		assert.Equal(t, []string{"one", "two"}, result.Values("X-Test"))
	})

	t.Run("rejects duplicate headers when disabled", func(t *testing.T) {
		t.Parallel()

		_, err := ParseHeaders([]string{
			"X-Test=one",
			"X-Test=two",
		}, false)

		assert.EqualError(t, err, "duplicate header key found: X-Test")
	})

	t.Run("duplicate detection is case insensitive", func(t *testing.T) {
		t.Parallel()

		_, err := ParseHeaders([]string{
			"x-test=one",
			"X-Test=two",
		}, false)

		assert.EqualError(t, err, "duplicate header key found: X-Test")
	})
}

func TestParseStatusCodes(t *testing.T) {
	t.Parallel()

	t.Run("valid status code", func(t *testing.T) {
		t.Parallel()

		statuses, err := ParseStatusCodes("200")

		require.NoError(t, err)
		assert.Equal(t, []int{200}, statuses)
	})

	t.Run("valid multiple status codes", func(t *testing.T) {
		t.Parallel()

		statuses, err := ParseStatusCodes("200,404,500")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 404, 500}, statuses)
	})

	t.Run("valid status code range", func(t *testing.T) {
		t.Parallel()

		statuses, err := ParseStatusCodes("200-202")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 201, 202}, statuses)
	})

	t.Run("valid multiple status code ranges", func(t *testing.T) {
		t.Parallel()

		statuses, err := ParseStatusCodes("200-202,300-301,500")

		require.NoError(t, err)
		assert.Equal(t, []int{200, 201, 202, 300, 301, 500}, statuses)
	})

	t.Run("invalid status code", func(t *testing.T) {
		t.Parallel()

		_, err := ParseStatusCodes("abc")

		assert.EqualError(t, err, "invalid status code: abc")
	})

	t.Run("invalid status range double dash", func(t *testing.T) {
		t.Parallel()

		_, err := ParseStatusCodes("200--202")

		assert.EqualError(t, err, "invalid status range: 200--202")
	})

	t.Run("invalid status range start greater than end", func(t *testing.T) {
		t.Parallel()

		_, err := ParseStatusCodes("201-200")

		assert.EqualError(t, err, "invalid status range: 201-200")
	})
}
