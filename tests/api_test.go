package tests

import (
	"net/http/httptest"
	"pharmacy-counter/internal/api"
	"testing"
)

func TestAPIIndex(t *testing.T) {
	p := makePharmacy(t)
	r := httptest.NewRecorder()
	api.New(p).Handler().ServeHTTP(r, httptest.NewRequest("GET", "/", nil))
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
}
