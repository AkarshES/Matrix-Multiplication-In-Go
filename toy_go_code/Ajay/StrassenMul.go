package main

import (
	"flag"
	"fmt"
	"time"
	"os"
)

const GRAIN int = 1024*1024 /* size of product(of dimensions) below which matmultleaf is used */
var mat1 string
var mat2 string

func init() {
	flag.StringVar(&mat1,"mat1","./data1.csv","Path to the CSV data file.")
	flag.StringVar(&mat2,"mat2","./data2.csv","Path to the CSV data file.")
}

func seqMatMult(m int, n int, p int, A [][]int, B [][]int, C [][]int) {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			C[i][j] = 0.0
			for k := 0; k < p; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
}

func matmultleaf(mf, ml, nf, nl, pf, pl int, A, B, C [][]int) {
	/* 
		subroutine that uses the simple triple loop to multiply 
		a submatrix from A with a submatrix from B and store the 
		result in a submatrix of C. 
	*/
	// mf, ml; /* first and last+1 i index */ 
	// nf, nl; /* first and last+1 j index */ 
	// pf, pl; /* first and last+1 k index */ 
	for i := mf; i < ml; i++ {
		for j := nf; j < nl; j++ {
			for k := pf; k < pl; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
}

func copyQtrMatrix(X [][]int, m int, Y [][]int, mf, nf int) {
	for i := 0; i < m; i++ {
		X[i] = Y[mf+i][nf:]
	}
}

func AddMats(T [][]int, m, n int, X [][]int, Y [][]int) {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			T[i][j] = X[i][j] + Y[i][j]
		}
	}
}

func SubMats(T [][]int, m, n int, X [][]int, Y [][]int) {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			T[i][j] = X[i][j] - Y[i][j]
		}
	}
}
func Allocate2DArray(m, n int) [][]int {
	temp := make([][]int, m)
	for i := 0; i < len(temp); i++ {
		temp[i] = make([]int, n)
	}
	return temp
}

func strassenMMult(mf, ml, nf, nl, pf, pl int, A, B, C [][]int) {
	if (ml-mf)*(nl-nf)*(pl-pf) < GRAIN {
		matmultleaf(mf, ml, nf, nl, pf, pl, A, B, C)
	} else {
		m2 := (ml - mf) / 2
		n2 := (nl - nf) / 2
		p2 := (pl - pf) / 2

		M1 := Allocate2DArray(m2, n2)
		M2 := Allocate2DArray(m2, n2)
		M3 := Allocate2DArray(m2, n2)
		M4 := Allocate2DArray(m2, n2)
		M5 := Allocate2DArray(m2, n2)
		M6 := Allocate2DArray(m2, n2)
		M7 := Allocate2DArray(m2, n2)

		A11 := make([][]int, m2)
		A12 := make([][]int, m2)
		A21 := make([][]int, m2)
		A22 := make([][]int, m2)

		B11 := make([][]int, p2)
		B12 := make([][]int, p2)
		B21 := make([][]int, p2)
		B22 := make([][]int, p2)

		C11 := make([][]int, m2)
		C12 := make([][]int, m2)
		C21 := make([][]int, m2)
		C22 := make([][]int, m2)

		tAM1 := Allocate2DArray(m2, p2)
		tBM1 := Allocate2DArray(p2, n2)
		tAM2 := Allocate2DArray(m2, p2)
		tBM3 := Allocate2DArray(p2, n2)
		tBM4 := Allocate2DArray(p2, n2)
		tAM5 := Allocate2DArray(m2, p2)
		tAM6 := Allocate2DArray(m2, p2)
		tBM6 := Allocate2DArray(p2, n2)
		tAM7 := Allocate2DArray(m2, p2)
		tBM7 := Allocate2DArray(p2, n2)

		copyQtrMatrix(A11, m2, A, mf, pf)
		copyQtrMatrix(A12, m2, A, mf, p2)
		copyQtrMatrix(A21, m2, A, m2, pf)
		copyQtrMatrix(A22, m2, A, m2, p2)

		copyQtrMatrix(B11, p2, B, pf, nf)
		copyQtrMatrix(B12, p2, B, pf, n2)
		copyQtrMatrix(B21, p2, B, p2, nf)
		copyQtrMatrix(B22, p2, B, p2, n2)

		copyQtrMatrix(C11, m2, C, mf, nf)
		copyQtrMatrix(C12, m2, C, mf, n2)
		copyQtrMatrix(C21, m2, C, m2, nf)
		copyQtrMatrix(C22, m2, C, m2, n2)

		// M1 = (A11 + A22)*(B11 + B22) 
		AddMats(tAM1, m2, p2, A11, A22)
		AddMats(tBM1, p2, n2, B11, B22)
		strassenMMult(0, m2, 0, n2, 0, p2, tAM1, tBM1, M1)

		//M2 = (A21 + A22)*B11 
		AddMats(tAM2, m2, p2, A21, A22)
		strassenMMult(0, m2, 0, n2, 0, p2, tAM2, B11, M2)

		//M3 = A11*(B12 - B22) 
		SubMats(tBM3, p2, n2, B12, B22)
		strassenMMult(0, m2, 0, n2, 0, p2, A11, tBM3, M3)

		//M4 = A22*(B21 - B11) 
		SubMats(tBM4, p2, n2, B21, B11)
		strassenMMult(0, m2, 0, n2, 0, p2, A22, tBM4, M4)

		//M5 = (A11 + A12)*B22 
		AddMats(tAM5, m2, p2, A11, A12)
		strassenMMult(0, m2, 0, n2, 0, p2, tAM5, B22, M5)

		//M6 = (A21 - A11)*(B11 + B12) 
		SubMats(tAM6, m2, p2, A21, A11)
		AddMats(tBM6, p2, n2, B11, B12)
		strassenMMult(0, m2, 0, n2, 0, p2, tAM6, tBM6, M6)

		//M7 = (A12 - A22)*(B21 + B22) 
		SubMats(tAM7, m2, p2, A12, A22)
		AddMats(tBM7, p2, n2, B21, B22)
		strassenMMult(0, m2, 0, n2, 0, p2, tAM7, tBM7, M7)

		for i := 0; i < m2; i++ {
			for j := 0; j < n2; j++ {
				C11[i][j] = M1[i][j] + M4[i][j] - M5[i][j] + M7[i][j]
				C12[i][j] = M3[i][j] + M5[i][j]
				C21[i][j] = M2[i][j] + M4[i][j]
				C22[i][j] = M1[i][j] - M2[i][j] + M3[i][j] + M6[i][j]
			}
		}
	}
}

func matmultS(m, n, p int, A, B, C [][]int) {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			C[i][j] = 0
		}
	}
	strassenMMult(0, m, 0, n, 0, p, A, B, C)
}

func CheckResults(m, n int, C, C1 [][]int) bool {
	// 
	// May need to take into consideration the floating point roundoff error 
	// due to parallel execution 
	// 
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if C[i][j] != C1[i][j] {
				fmt.Printf("C is\n")
				for i := 0; i < m; i++ {
					for j := 0; j < n; j++ {
						fmt.Printf("%v ", C[i][j])
					}
					fmt.Printf("\n")
				}
				fmt.Printf("C1 is\n")
				for i := 0; i < m; i++ {
					for j := 0; j < n; j++ {
						fmt.Printf("%v ", C1[i][j])
					}
					fmt.Printf("\n")
				}
				return true
			}
		}
	}
	return false
}

func main() {
	N := 0
	flag.Parse()
	A := OpenCsv(mat1)
	B := OpenCsv(mat2)
	M := A.Rows
	P := B.Columns
	if(A.Columns==B.Rows){
		N = A.Columns
	} else {
		fmt.Println("These matrices cannot be multiplied, %s has %d columns and %s has %d rows",mat1,A.Columns,mat2,B.Rows)
		os.Exit(1)
	}
//	A := Allocate2DArray(M, P)
//	B := Allocate2DArray(P, N)
	C := Allocate2DArray(M, N)
	C4 := Allocate2DArray(M, N)

//	for i := 0; i < M; i++ {
//		for j := 0; j < P; j++ {
//			A[i][j] = 5.0 - ((rand.Int() % 100) / 10.0)
//		}
//	}

//	for i := 0; i < P; i++ {
//		for j := 0; j < N; j++ {
//			if i == j {
//				B[i][j] = 1
//			} else {
//				B[i][j] = 0
//			}
//			//B[i][j] = 5.0 - ((rand.Int() % 100) / 10.0)
//		}
//	}

	fmt.Printf("Execute Standard matmult\n\n")
	before := time.Now()
	seqMatMult(M, N, P, A.Data, B.Data, C)
	after := time.Now()
	fmt.Printf("Standard matrix function done in %v s\n\n\n", after.Sub(before).Seconds())

	before = time.Now()
	matmultS(M, N, P, A.Data, B.Data, C4)
	after = time.Now()
	fmt.Printf("Strassen matrix function done in %v s\n\n\n", after.Sub(before).Seconds())

	if CheckResults(M, N, C, C4) {
		fmt.Printf("Error in matmultS\n\n")
	} else {
		fmt.Printf("OKAY\n\n")
	}
}
