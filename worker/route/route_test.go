package route

import "testing"

func TestFindPath(t *testing.T) {
	path := FindPath("https://en.wikipedia.org/wiki/Apple_A4",
		"https://en.wikipedia.org/wiki/Acorn_Computers")

	for i := len(path) - 1; i >= 0; i-- {
		t.Log(path[i])
	}

}
