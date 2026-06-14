package tigergraphx

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gospacex/graphx"
)

func TestBuild_TigerAddressRequired(t *testing.T) {
	_, err := Build(graphx.Config{})
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestBuild_TigerDefaultScheme(t *testing.T) {
	tcfg, err := Build(graphx.Config{Address: "localhost:14240"})
	if err != nil {
		t.Fatal(err)
	}
	expected := "http://127.0.0.1:14240"
	if tcfg.BaseURL != expected {
		t.Fatalf("expected %s, got %s", expected, tcfg.BaseURL)
	}
}

func TestBuild_TigerTLSScheme(t *testing.T) {
	tcfg, err := Build(graphx.Config{Address: "127.0.0.1:14240", TLS: true})
	if err != nil {
		t.Fatal(err)
	}
	if tcfg.BaseURL != "https://127.0.0.1:14240" {
		t.Fatalf("expected https://127.0.0.1:14240, got %s", tcfg.BaseURL)
	}
}

func TestNew_TigerNilCtx(t *testing.T) {
	defer Reset()
	tg, err := New(nil, graphx.Config{Address: "127.0.0.1:14240"})
	if err != nil {
		t.Fatal("tigergraph client creation should not fail without live connection", err)
	}
	if tg.HTTPClient() == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestGet_TigerNamedSingleton(t *testing.T) {
	defer Reset()
	tg, err := Get(context.Background(), graphx.Config{Name: "t1", Address: "127.0.0.1:14240"})
	if err != nil {
		t.Fatal("tigergraph singleton should not fail without live connection", err)
	}
	if tg.HTTPClient() == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestReset_Tiger(t *testing.T) {
	Reset()
}

// TestBuildLoginBody_QuoteEscape verifies that login credentials containing
// JSON-significant characters (quotes, backslashes, control characters) are
// properly escaped. Regression test for the fmt.Sprintf JSON injection risk
// (close-issues F-1).
func TestBuildLoginBody_QuoteEscape(t *testing.T) {
	cases := []struct {
		name     string
		user     string
		password string
	}{
		{"plain", "alice", "secret"},
		{"double-quote", `evil"name`, `pass"word`},
		{"backslash", `domain\user`, `back\slash`},
		{"newline", "alice", "line1\nline2"},
		{"unicode", "用户", "密码"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := buildLoginBody(tc.user, tc.password)
			if err != nil {
				t.Fatalf("buildLoginBody: %v", err)
			}
			var got struct {
				User     string `json:"user"`
				Password string `json:"password"`
			}
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("body is not valid JSON: %v\nbody: %s", err, body)
			}
			if got.User != tc.user {
				t.Errorf("user roundtrip: want %q, got %q", tc.user, got.User)
			}
			if got.Password != tc.password {
				t.Errorf("password roundtrip: want %q, got %q", tc.password, got.Password)
			}
		})
	}
}

// TestBuildLoginBody_NoRawSprintf verifies the body never contains the
// raw user/password string when those values include double-quotes (which
// would indicate the JSON was assembled via fmt.Sprintf concatenation).
func TestBuildLoginBody_NoRawSprintf(t *testing.T) {
	body, err := buildLoginBody(`u"ser`, `p"ass`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), `p"ass`) {
		t.Errorf("body should JSON-escape quote, got: %s", body)
	}
}
