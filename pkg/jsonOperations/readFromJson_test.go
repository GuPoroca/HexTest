package jsonOperations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/GuPoroca/HexTest/pkg/typeDefines"
)

func TestReadJSONParsesProject(t *testing.T) {
	p := ReadJSON(filepath.Join("testdata", "sample-project.json"))

	if p.Name != "Sample Project" {
		t.Fatalf("Name: got %q, want %q", p.Name, "Sample Project")
	}
	if p.Url != "http://localhost:3443" {
		t.Fatalf("Url: got %q", p.Url)
	}
	if len(p.Suites) != 1 || len(p.Suites[0].Tests) != 1 {
		t.Fatalf("expected 1 suite with 1 test, got %d suites", len(p.Suites))
	}

	test := p.Suites[0].Tests[0]
	if test.Method != "GET" || test.Api_endpoint != "/test/schema" {
		t.Fatalf("unexpected test fields: %+v", test)
	}
	if len(test.Asserts) != 2 {
		t.Fatalf("expected 2 asserts, got %d", len(test.Asserts))
	}
	if got := test.Asserts[1].Checks[0].Expected; len(got) != 3 {
		t.Fatalf("expected 3 expected values on the containsKey check, got %v", got)
	}
}

func TestWriteJSONRoundTrips(t *testing.T) {
	// WriteJSON writes "<Project.Name>.json" into the current directory, so run
	// it from a throwaway directory.
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	want := typeDefines.Project{
		Name: "RoundTrip",
		Url:  "http://example.test",
		Suites: []typeDefines.Suite{
			{Name: "S1", Tests: []typeDefines.Test{{Name: "T1", Method: "GET", Api_endpoint: "/"}}},
		},
	}
	WriteJSON("ignored", want)

	got := ReadJSON(filepath.Join(dir, "RoundTrip.json"))
	if got.Name != want.Name || got.Url != want.Url {
		t.Fatalf("round trip mismatch: got %+v", got)
	}
	if len(got.Suites) != 1 || got.Suites[0].Tests[0].Name != "T1" {
		t.Fatalf("round trip lost nested data: %+v", got.Suites)
	}
}
