package result

import (
	"errors"
	"testing"
)

func TestValidationResult(t *testing.T) {
	r := &ValidationResult{}
	if r.HasErrors() {
		t.Fatal("new result should report no errors")
	}

	r.AddSuccess("ok")
	r.AddInfo("fyi")
	if r.HasErrors() {
		t.Fatal("success/info must not count as errors")
	}

	r.AddFailure(errors.New("boom"))
	if !r.HasErrors() {
		t.Fatal("expected HasErrors after AddFailure")
	}

	if len(r.Success) != 1 || len(r.Info) != 1 || len(r.Failure) != 1 {
		t.Fatalf("unexpected counts: success=%d info=%d failure=%d",
			len(r.Success), len(r.Info), len(r.Failure))
	}
}
