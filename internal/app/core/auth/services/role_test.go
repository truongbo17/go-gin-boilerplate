package services

import "testing"

func TestReservedAdminSlug(t *testing.T) {
	for _, slug := range []string{"admin", "ADMIN", " admin "} {
		if !isReservedAdminSlug(slug) {
			t.Errorf("admin slug %q was not reserved", slug)
		}
	}
	if isReservedAdminSlug("editor") {
		t.Fatal("ordinary role was reserved")
	}
}
