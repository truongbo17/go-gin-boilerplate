package requests

import "testing"

func TestRoleSlugValidation(t *testing.T) {
	for _, slug := range []string{"ADMIN", "admin ", "bad slug", ""} {
		request := CreateRoleRequest{Name: "Example", Slug: slug}
		if err := request.Validate("en"); err == nil {
			t.Errorf("invalid slug %q was accepted", slug)
		}
	}
	request := CreateRoleRequest{Name: "Example", Slug: "team_admin"}
	if err := request.Validate("en"); err != nil {
		t.Fatal(err)
	}
}
