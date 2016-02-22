package gitter

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNew_setToken(t *testing.T) {
	g := New("abc")
	if g.config.token != "abc" {
		t.Logf("Expected %v, got %v", "abc", g.config.token)
	}
}

func TestNew_setAPIBaseURL(t *testing.T) {
	g := New("abc")
	if g.config.apiBaseURL != apiBaseURL {
		t.Logf("Expected %v, got %v", apiBaseURL, g.config.apiBaseURL)
	}
}

func TestNew_setStreamBaseURL(t *testing.T) {
	g := New("abc")
	if g.config.streamBaseURL != streamBaseURL {
		t.Logf("Expected %v, got %v", streamBaseURL, g.config.streamBaseURL)
	}
}

func TestGitter_SetClient(t *testing.T) {
	setup()
	defer teardown()

	c := &http.Client{}
	gitter.SetClient(c)

	if gitter.config.client != c {
		t.Logf("Expected %v, got %v", c, gitter.config.client)
	}
}

func TestGetUser_userData(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
            [
              {
                "id": "123",
                "username": "fooBar"
              }
            ]
        `)
	})

	u, err := gitter.GetUser()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if u.ID != "123" {
		t.Errorf("Expected %v, got %v", "123", u.ID)
	}

	if u.Username != "fooBar" {
		t.Errorf("Expected %v, got %v", "fooBar", u.Username)
	}
}

func TestGetUser_apiEmptyResponseError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "[]")
	})

	u, err := gitter.GetUser()
	if u != nil {
		t.Errorf("Expected %v, got %v", nil, u)
	}

	if err.Error() != "Failed to retrieve current user" {
		t.Errorf("Expected %v, got %v", "Failed to retrieve current user", err)
	}
}

func TestGetUser_apiError(t *testing.T) {
	setup()
	defer teardown()

	wanted := "json: cannot unmarshal object into Go value of type []gitter.User"

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
            {
              "error": "Unauthorized"
            }
        `)
	})

	u, err := gitter.GetUser()
	if u != nil {
		t.Errorf("Expected %v, got %v", nil, u)
	}

	if err.Error() != wanted {
		t.Errorf("Expected %v, got %v", wanted, err)
	}
}

func TestGetUserRooms(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/abc/rooms/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
            [
                {
                    "id": "xyz"
                },
                {
                    "id": "cde"
                }
            ]

            `)
	})

	r, err := gitter.GetUserRooms("abc")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if len(r) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(r))
	}

	if r[0].ID != "xyz" {
		t.Errorf("Expected %v, got %v", "xyz", r[0].ID)
	}

	if r[1].ID != "cde" {
		t.Errorf("Expected %v, got %v", "cde", r[1].ID)
	}
}

func TestGetRooms(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
            [
                {
                    "id": "xyz"
                }
            ]

            `)
	})

	r, err := gitter.GetRooms()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if len(r) != 1 {
		t.Errorf("Expected %v, got %v", 2, len(r))
	}

	if r[0].ID != "xyz" {
		t.Errorf("Expected %v, got %v", "xyz", r[0].ID)
	}
}

func TestGetRoom(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rooms/xyz/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
            {
                "id": "xyz"
            }
            `)
	})

	r, err := gitter.GetRoom("xyz")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if r.ID != "xyz" {
		t.Errorf("Expected %v, got %v", "xyz", r.ID)
	}
}
