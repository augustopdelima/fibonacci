package main

import (
	bottomup "fibonacci/bottom-up"
	"fibonacci/recursiva"
	topdown "fibonacci/top-down"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type AlgoritmoFibonnaci func(int) int

type ResultadoLinha struct {
	Tamanho int
	Tempos  map[string]float64
}

func testarTempo(algoritmo AlgoritmoFibonnaci, numero int) float64 {

	tempoInicio := time.Now()

	_ = algoritmo(numero)

	tempoFim := time.Since(tempoInicio)

	tempoMiliSegundos := float64(tempoFim.Nanoseconds()) / 1e6

	return tempoMiliSegundos

}

func criarArquivoExcel(algoritmos []string, resultados []ResultadoLinha) {
	arquivoExcel := excelize.NewFile()
	defer func() {
		if err := arquivoExcel.Close(); err != nil {
			fmt.Println("Erro ao fechar arquivo:", err)
		}
	}()

	planilha := "Sheet1"

	cabecalho := []string{"Tamanho"}
	cabecalho = append(cabecalho, algoritmos...)

	if err := arquivoExcel.SetSheetRow(planilha, "A1", &cabecalho); err != nil {
		fmt.Println("Erro ao escrever cabeçalho:", err)
	}

	linhaExcel := 2

	for _, linha := range resultados {
		cell, err := excelize.CoordinatesToCellName(1, linhaExcel)

		valores := make([]any, 0, len(algoritmos)+1)
		valores = append(valores, linha.Tamanho)

		for _, nomeAlg := range algoritmos {
			valores = append(valores, linha.Tempos[nomeAlg])
		}

		if err != nil {
			fmt.Println("Erro ao calcular célula:", err)
			return
		}

		if err := arquivoExcel.SetSheetRow(planilha, cell, &valores); err != nil {
			fmt.Println("Erro ao escrever linha:", err)
			return
		}

		linhaExcel++
	}

	criarGrafico(algoritmos, arquivoExcel, planilha, linhaExcel)
}

func criarGrafico(algoritmos []string, arquivoExcel *excelize.File, planilha string, linhaExcel int) {
	var series []excelize.ChartSeries

	for i := range algoritmos {
		letraColuna, _ := excelize.ColumnNumberToName(i + 2)
		serie := excelize.ChartSeries{
			Name:       fmt.Sprintf("%s!$%s$1", planilha, letraColuna),
			Categories: fmt.Sprintf("%s!$A$2:$A$%d", planilha, linhaExcel-1),
			Values:     fmt.Sprintf("%s!$%s$2:$%s$%d", planilha, letraColuna, letraColuna, linhaExcel-1),
		}
		series = append(series, serie)
	}

	err := arquivoExcel.AddChart(planilha, "J2", &excelize.Chart{
		Type:   excelize.Line,
		Series: series,
		Title: []excelize.RichTextRun{
			{Text: "Tempo x Tamanho Elementos"},
		},
		XAxis: excelize.ChartAxis{Title: []excelize.RichTextRun{
			{Text: "Tamanho"},
		}},
		YAxis: excelize.ChartAxis{Title: []excelize.RichTextRun{
			{Text: "Tempo em (ms)"},
		}},
	})

	if err != nil {
		fmt.Println("Erro ao adicionar gráfico:", err)
		return
	}

	if err := arquivoExcel.SaveAs("resultado.xlsx"); err != nil {
		fmt.Println("Erro ao salvar arquivo:", err)
		return
	}

	fmt.Println("Arquivo resultado.xlsx gerado com gráfico!")
}

const TAMANHO_MAXIMO = 40
const TAMANHO_MINIMO = 0

func main() {

	algoritmos := map[string]AlgoritmoFibonnaci{
		"Recursiva": recursiva.FibonacciRecursiva,
		"Bottom-Up": bottomup.FibonacciBottomUp,
		"Top-Down":  topdown.FibonacciTopDown,
	}

	ordemAlgoritmos := []string{
		"Recursiva",
		"Bottom-Up",
		"Top-Down",
	}

	resultados := make([]ResultadoLinha, TAMANHO_MAXIMO-TAMANHO_MINIMO+1)

	for _, nome := range ordemAlgoritmos {

		algoritmo := algoritmos[nome]

		for i := TAMANHO_MINIMO; i <= TAMANHO_MAXIMO; i++ {
			indice := i - TAMANHO_MINIMO

			if resultados[indice].Tempos == nil {
				resultados[indice] = ResultadoLinha{
					Tamanho: i,
					Tempos:  make(map[string]float64),
				}
			}

			tempo := testarTempo(algoritmo, i)

			resultados[indice].Tempos[nome] = tempo
		}
	}

	criarArquivoExcel(ordemAlgoritmos, resultados)
}
