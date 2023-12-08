package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

func main() {
	perCapita, err := ReadEmissionsCsv(filepath.Join("emissions", "co-emissions-per-capita.csv"))
	cumulative, err := ReadEmissionsCsv(filepath.Join("emissions", "cumulative-co-emissions.csv"))
	combined := CombineCSVs(perCapita, cumulative, 3)

	f, err := os.Create(filepath.Join("emissions", "per_capita_cumulative.csv"))
	Must(err)

	w := csv.NewWriter(f)
	Must(w.WriteAll(combined))
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadEmissionsCsv(path string) ([][]string, error) {
	f, err := os.Open(path)

	Must(err)

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	Must(err)
	return rows, nil
}

func CombineCSVs(capita [][]string, cumulative [][]string, appendIdx int) [][]string {
	capitaValues := map[[3]string]string{}
	cumuValues := map[[3]string]string{}
	for i := 0; i < len(capita); i++ {
		capitaValues[[3]string{capita[i][0], capita[i][1], capita[i][2]}] = capita[i][3]
	}
	for i := 0; i < len(cumulative); i++ {
		cumuValues[[3]string{cumulative[i][0], cumulative[i][1], cumulative[i][2]}] = cumulative[i][3]
	}

	// adjust header
	results := make([][]string, 1, len(capita))
	results[0] = append(capita[0], cumulative[0][3])

	for i := 1; i < len(capita); i++ {
		if val, ok := cumuValues[[3]string{capita[i][0], capita[i][1], capita[i][2]}]; ok {
			results = append(results, append(capita[i], val))
		}
	}
	return results
}
