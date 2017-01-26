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

func TestGetMessages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rooms/xyz/chatMessages", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
            [{
                "id": "666"
            }]
        `)
	})

	m, err := gitter.GetMessages("xyz", nil)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if len(m) != 1 {
		t.Errorf("Expected %v, got %v", 1, len(m))
	}

	if m[0].ID != "666" {
		t.Errorf("Expected %v, got %v", "666", m[0].ID)
	}
}

func TestGetMessages_limit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rooms/xyz/chatMessages", func(w http.ResponseWriter, r *http.Request) {
		if ok := len(r.URL.Query().Get("limit")) > 0; !ok {
			t.Errorf("Expected %v, got %v", 1, 0)
		}
		fmt.Fprint(w, `
            [{
                "id": "666"
            },
            {
                "id": "112"
            }]
        `)
	})

	p := &Pagination{
		Limit: 2,
	}

	m, err := gitter.GetMessages("xyz", p)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if len(m) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(m))
	}
}

func TestGetMessage(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rooms/xyz/chatMessages/666", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
            {
                "id": "666"
            }
        `)
	})

	m, err := gitter.GetMessage("xyz", "666")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if m.ID != "666" {
		t.Errorf("Expected %v, got %v", "666", m.ID)
	}
}

func TestSendMessage(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/rooms/xyz/chatMessages", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := gitter.SendMessage("xyz", "test message.")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
}

func TestGetResponse(t *testing.T) {
	setup()
	defer teardown()

	r, err := gitter.getResponse(gitter.config.apiBaseURL, nil)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	if r.Request.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected %v, got %v", "application/json", r.Request.Header.Get("Content-Type"))
	}

	if r.Request.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected %v, got %v", "application/json", r.Request.Header.Get("Accept"))
	}

	if r.Request.Header.Get("Authorization") != "Bearer abc" {
		t.Errorf("Expected %v, got %v", "Bearer abc", r.Request.Header.Get("Authorization"))
	}

	if r.Request.URL.String() != gitter.config.apiBaseURL {
		t.Errorf("Expected %v, got %v", gitter.config.apiBaseURL, r.Request.URL.String())
	}
}

func TestGet_endpoint(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	_, err := gitter.get(gitter.config.apiBaseURL)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
}

func TestGet_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})

	b, _ := gitter.get(gitter.config.apiBaseURL)
	if b != nil {
		t.Errorf("Expected %v, got %v", nil, b)
	}
}

func TestPost(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected %v, got %v", "application/json", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected %v, got %v", "application/json", r.Header.Get("Accept"))
		}

		if r.Header.Get("Authorization") != "Bearer abc" {
			t.Errorf("Expected %v, got %v", "Bearer abc", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
	})

	_, err := gitter.post(gitter.config.apiBaseURL, []byte{})
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
}
