package download

import "testing"

func TestParseHTTP(t *testing.T) {
	in, err := Parse("https://example.com/file.zip")
	if err != nil {
		t.Fatal(err)
	}
	if in.Kind != KindHTTP {
		t.Fatalf("got %v", in.Kind)
	}
}

func TestParseMagnet(t *testing.T) {
	in, err := Parse("magnet:?xt=urn:btih:abc")
	if err != nil {
		t.Fatal(err)
	}
	if in.Kind != KindMagnet {
		t.Fatalf("got %v", in.Kind)
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := Parse("  ")
	if err != ErrEmpty {
		t.Fatalf("got %v", err)
	}
}

func TestParseUnrecognized(t *testing.T) {
	_, err := Parse("not-a-url")
	if err != ErrUnrecognized {
		t.Fatalf("got %v", err)
	}
}
