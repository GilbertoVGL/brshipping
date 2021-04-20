package rate

import (
	"testing"
)

type requestTest struct {
	contract, password, service                    string
	notify, inHands                                bool
	origin, destiny, weight, shape                 uint64
	length, height, width, diameter, declaredValue float64
}

func TestPrecoPrazo(t *testing.T) {
	var tests = []struct {
		request requestTest
		want    Response
	}{
		{
			requestTest{
				contract:      "",
				password:      "",
				service:       "04014",
				notify:        false,
				inHands:       false,
				origin:        74110100,
				destiny:       75113180,
				weight:        2,
				shape:         1,
				length:        15,
				height:        15,
				width:         15,
				diameter:      0,
				declaredValue: 0,
			},
			Response{
				Code:                 4014,
				Value:                24.90,
				Deadline:             1,
				DeliverInHandsValue:  0,
				NotifyReceivingValue: 0,
				DeclaredValue:        0,
				HomeDeliver:          "S",
				SaturdayDeliver:      "S",
				Error:                0,
				ErrorMessage:         "",
				PlainValue:           24.90,
			},
		},
		{
			requestTest{
				contract:      "",
				password:      "",
				service:       "service",
				notify:        false,
				inHands:       false,
				origin:        74110100,
				destiny:       75113180,
				weight:        2,
				shape:         15,
				length:        15,
				height:        15,
				width:         15,
				diameter:      0,
				declaredValue: 0,
			},
			Response{
				Code:                 0,
				Value:                0,
				Deadline:             0,
				DeliverInHandsValue:  0,
				NotifyReceivingValue: 0,
				DeclaredValue:        0,
				HomeDeliver:          "",
				SaturdayDeliver:      "",
				Error:                99,
				ErrorMessage:         "Erro inesperado. Descrição: Input string was not in a correct format.",
				PlainValue:           0,
			},
		},
	}

	for _, tt := range tests {
		r := tt.request
		ans, err := Calc(r.contract, r.password, r.service, r.notify,
			r.inHands, r.origin, r.destiny, r.weight, r.shape, r.length,
			r.height, r.width, r.diameter, r.declaredValue)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if ans != tt.want {
			t.Errorf("got: %v\nwant: %v\n", ans, tt.want)
		}
	}
}
