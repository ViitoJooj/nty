package services

import "testing"

func TestSemverRe(t *testing.T) {
	good := []string{"v0.0.0", "v1.2.3", "v10.20.30"}
	bad := []string{"1.2.3", "v1.2", "version v1.2.3", "v1.2.3-rc1", ""}

	for _, s := range good {
		if !semverRe.MatchString(s) {
			t.Errorf("esperava válido: %q", s)
		}
	}
	for _, s := range bad {
		if semverRe.MatchString(s) {
			t.Errorf("esperava inválido: %q", s)
		}
	}
}
