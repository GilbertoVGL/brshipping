package precoprazo

import (
	"testing"
)

func TestPrecoPrazo(t *testing.T) {
	var tests = []struct {
		request request
		want    Response
	}{
		{
			New("", "", "04014", false, false, 74110100, 75113180, 2, 1, 15, 15, 15, 0, 0),
			Response{
				Codigo:                4014,
				Valor:                 24.90,
				PrazoEntrega:          1,
				ValorMaoPropria:       0,
				ValorAvisoRecebimento: 0,
				ValorValorDeclarado:   0,
				EntregaDomiciliar:     "S",
				EntregaSabado:         "S",
				Erro:                  0,
				MsgErro:               "",
				ValorSemAdicionais:    24.90,
			},
		},
		{
			New("", "", "service", false, false, 74110100, 75113180, 2, 1, 15, 15, 15, 0, 0),
			Response{
				Codigo:                0,
				Valor:                 0,
				PrazoEntrega:          0,
				ValorMaoPropria:       0,
				ValorAvisoRecebimento: 0,
				ValorValorDeclarado:   0,
				EntregaDomiciliar:     "",
				EntregaSabado:         "",
				Erro:                  99,
				MsgErro:               "Erro inesperado. Descrição: Input string was not in a correct format.",
				ValorSemAdicionais:    0,
			},
		},
	}

	// New("", "", "", false, false, 0, 1, 1, 1, 1.0, 1.0, 1.0, 1.0, 1.0)

	for _, tt := range tests {
		ans, err := Calc(tt.request)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if ans != tt.want {
			t.Errorf("got: %v\nwant: %v\n", ans, tt.want)
		}
	}
}
