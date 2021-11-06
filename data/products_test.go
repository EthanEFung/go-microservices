package data

import "testing"

func TestValidation(t *testing.T) {
	p := &Product{
		Name:  "Frappe",
		Price: 1,
		SKU:   "abs-b-cs",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
