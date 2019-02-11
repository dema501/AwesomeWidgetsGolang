package types

import (
	"testing"
)

func TestProduce(t *testing.T) {
	dataCh := make(chan Widget, 8)

	p := NewProducer(1)
	p.Produce(4, 4, dataCh)

	close(dataCh)

	if len(dataCh) != 4 {
		t.Error("expected 4 widgets;", len(dataCh))
	}

	for elem := range dataCh {
		if elem.broken != true {
			t.Error("expected broken widget;", elem)
		}
	}
}
