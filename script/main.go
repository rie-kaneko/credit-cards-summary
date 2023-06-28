package main

import (
	"encoding/csv"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"os"
	"time"
)

func main() {
	for i := 0; i < 20; i++ {
		generateData(gofakeit.UUID())
	}
}

func generateData(name string) {
	for _, t := range []string{"credit", "debit"} {
		data := make([][]string, 0, 20)
		data = append(data, []string{"id", "date", "transaction"})

		for i := 0; i < gofakeit.IntRange(10, 100); i++ {
			id := fmt.Sprintf("%d", i)
			m := gofakeit.Month()
			var d int

			switch m {
			case int(time.February):
				d = gofakeit.IntRange(1, 27)
			default:
				d = gofakeit.IntRange(1, 30)
			}

			date := fmt.Sprintf("%d/%d/%d", m, d, time.Now().Year())
			trans := fmt.Sprintf("%.1f", gofakeit.Float64Range(-100, 100))

			data = append(data, []string{id, date, trans})
		}

		file, err := os.Create(fmt.Sprintf("resources/to_process/%s_%s.csv", name, t))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		for _, row := range data {
			err = writer.Write(row)
			if err != nil {
				panic(err)
			}
		}
	}
}
