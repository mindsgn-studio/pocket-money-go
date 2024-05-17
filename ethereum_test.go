package ethereum

import "testing"

func TestCreateNewWallet(t *testing.T) {
	want := ""
	got, err := CreateNewWallet()
	if err != nil {
		t.Errorf("got %q, wanted %q", got, want)
	}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
