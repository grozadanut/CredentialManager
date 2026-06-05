package credential

import (
	"moqui/runtime/component/CredentialManager/src/test/go/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminCanReadAll(t *testing.T) {
	resp := NewRequest(t, "GET",
		"/rest/s1/CredentialManager/credential", "john.doe", "moqui", map[string]string{})

	assert.Equal(t, 3, len(resp.GetList("resultList")))
	for _, value := range resp.GetList("resultList") {
		if value.Get("name") == "user" {
			assert.Equal(t, value.Get("value"), "john.doe")
		} else if value.Get("name") == "password" {
			assert.Equal(t, value.Get("value"), "moqui")
		} else if value.Get("name") == "token" {
			assert.Equal(t, value.Get("value"), "123tokenstring456")
		}
	}
}

func TestReaderCanReadCred1(t *testing.T) {
	resp := NewRequest(t, "GET", "/rest/s1/CredentialManager/credential", "cred.reader1", "moqui", map[string]string{})
	resultList := resp.GetList("resultList")

	assert.Equal(t, 2, len(resultList))
	userField := utils.FindByKeyValue(resultList, "name", "user")
	if userField == nil {
		t.Fatal("user field not found")
	}
	passField := utils.FindByKeyValue(resultList, "name", "password")
	if passField == nil {
		t.Fatal("password field not found")
	}
	assert.Equal(t, "john.doe", userField.Get("value"))
	assert.Equal(t, "moqui", passField.Get("value"))
}

func TestReaderCannotReadCred2(t *testing.T) {
	resp := NewRequest(t, "GET", "/rest/s1/CredentialManager/credential", "cred.reader1", "moqui",
		map[string]string{"credentialId": "CRED2"})
	resultList := resp.GetList("resultList")

	assert.Equal(t, 0, len(resultList))
}

func TestReader2CanReadCred2(t *testing.T) {
	resp := NewRequest(t, "GET", "/rest/s1/CredentialManager/credential", "cred.reader2", "moqui", map[string]string{})
	resultList := resp.GetList("resultList")

	assert.Equal(t, 1, len(resultList))
	tokenField := utils.FindByKeyValue(resultList, "name", "token")
	if tokenField == nil {
		t.Fatal("token field not found")
	}
	assert.Equal(t, "123tokenstring456", tokenField.Get("value"))
}

func TestReader2CannotReadCred1(t *testing.T) {
	resp := NewRequest(t, "GET", "/rest/s1/CredentialManager/credential", "cred.reader2", "moqui",
		map[string]string{"credentialId": "CRED1"})
	resultList := resp.GetList("resultList")

	assert.Equal(t, 0, len(resultList))
}

func TestReaderCannotReadExpiredCred(t *testing.T) {
	resp := NewRequest(t, "GET", "/rest/s1/CredentialManager/credential", "cred.reader1", "moqui",
		map[string]string{"credentialId": "CRED1_EXPIRED"})
	resultList := resp.GetList("resultList")

	assert.Equal(t, 0, len(resultList))
}

func TestReaderCannotReadExpiredPermissionCred(t *testing.T) {
	resp := NewRequest(t, "GET", "/rest/s1/CredentialManager/credential", "cred.reader1", "moqui",
		map[string]string{"credentialId": "CRED1_EXPIRED_PERMISSION"})
	resultList := resp.GetList("resultList")

	assert.Equal(t, 0, len(resultList))
}

func TestIfDateIsOldReadOldField(t *testing.T) {
	resp := NewRequest(t, "GET", "/rest/s1/CredentialManager/credential", "cred.reader1", "moqui",
		map[string]string{"credentialId": "CRED1", "date": "2026-03-01T00:00:00Z"})
	resultList := resp.GetList("resultList")

	assert.Equal(t, 1, len(resultList))
	passField := utils.FindByKeyValue(resultList, "name", "password")
	if passField == nil {
		t.Fatal("password field not found")
	}
	assert.Equal(t, "oldpassword", passField.Get("value"))
}
