package auth

import (
	"testing"
)

func TestHashingAndComparePwd(t *testing.T) {
	cases := []struct {
		password string
		hash     string
	}{
		{
			password: "1234",
		},
		{
			password: "password",
		},
		{
			password: "Yhdjls9@aahsk\\!ahjdgzxJGSFz",
		},
	}

	for _, c := range cases {
		actual, err := HashPassword(c.password)
		if err != nil {
			t.Errorf("Error hashing password: %v", err)
		}

		if CheckPasswordHash(c.password, actual) != true {
			t.Errorf("Error checking password\n")
		}
	}
}
