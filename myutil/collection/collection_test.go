package collection

import "testing"

func TestIndex(t *testing.T) {
	t.Run("string", testIndexByString)
	t.Run("Int", testIndexByInt)
	t.Run("Struct", testIndexByStrcut)
}

func testIndexByString(t *testing.T) {
	expected := 1
	slc := []string{"a", "b", "c"}
	actual := Index(slc, "b")
	if expected != actual {
		t.Errorf("Test failed, expected:'%d', got: '%d'", expected, actual)
	}
}

func testIndexByInt(t *testing.T) {
	expected := 1
	slc := []int{1, 2, 3}
	actual := Index(slc, 2)
	if expected != actual {
		t.Errorf("Test failed, expected:'%d', got: '%d'", expected, actual)
	}
}

func testIndexByStrcut(t *testing.T) {
	expected := 1
	type user struct {
		Name string
	}
	slc := []user{
		user{"a"},
		user{"b"},
		user{"c"},
	}
	actual := Index(slc, user{"b"})
	if expected != actual {
		t.Errorf("Test failed, expected:'%d', got: '%d'", expected, actual)
	}
}

func TestIn(t *testing.T) {
	expected := true
	slc := []string{"a", "b", "c"}
	actual := In(slc, "b")
	if expected != actual {
		t.Errorf("Test failed, expected:'%t', got: '%t'", expected, actual)
	}
}
