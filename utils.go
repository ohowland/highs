package highs

type PackedMatrix struct {
	arStart []int
	arIndex []int
	arValue []float64
}

// packRows returns the rows in a flat packed form use in the HiGHS c interface
func packMatrix(matrix [][]float64) PackedMatrix {
	arStart := []int{}
	arIndex := []int{}
	arValue := []float64{}

	idx := 0
	for _, cols := range matrix {
		arStart = append(arStart, idx)
		for i, v := range cols {
			if v != 0 {
				arIndex = append(arIndex, i)
				arValue = append(arValue, v)
			}
		}
		idx = len(arIndex)
	}

	return PackedMatrix{arStart, arIndex, arValue}
}

// separateBounds splits a bounded rows [lb, row..., ub] into its components lb, rows, ub
func separateBounds(bound_rows [][]float64) ([][]float64, []float64, []float64) {
	col_size := len(bound_rows[0])

	rows := [][]float64{}
	lbs := []float64{}
	ubs := []float64{}
	for _, row := range bound_rows {
		rows = append(rows, row[1:col_size-1])
		lbs = append(lbs, row[0])
		ubs = append(ubs, row[col_size-1])
	}

	return rows, lbs, ubs
}
