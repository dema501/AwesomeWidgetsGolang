package producer

import (
	"testing"

	"github.com/dema501/AwesomeWidgetsGolang/internal/pkg/types"
)

func TestProduce(t *testing.T) {
	dataCh := make(chan *types.Widget, 8)

	p := NewProducer(1)
	p.Produce(4, 4, dataCh)

	close(dataCh)

	if len(dataCh) != 4 {
		t.Error("expected 4 widgets;", len(dataCh))
	}

	for elem := range dataCh {
		if elem.Broken != true {
			t.Error("expected broken widget;", elem)
		}
	}
}
