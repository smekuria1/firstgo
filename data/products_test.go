package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "yoo",
		Price: 12,
		SKU:   "abc-acd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
