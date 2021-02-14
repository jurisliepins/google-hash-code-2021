package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Problem struct {
	M  int
	T2 int
	T3 int
	T4 int
	P  [][]string
}

type Delivery struct {
	T int
	P []int
}
type Solution struct {
	D []Delivery
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func read(filename string) Problem {
	file, err := os.Open(filename)
	defer file.Close()
	must(err)

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	split := strings.Split(scanner.Text(), " ")

	M, err := strconv.Atoi(split[0])
	must(err)
	T2, err := strconv.Atoi(split[1])
	must(err)
	T3, err := strconv.Atoi(split[2])
	must(err)
	T4, err := strconv.Atoi(split[3])
	must(err)

	P := make([][]string, M)

	for i := 0; i < M; i++ {
		scanner.Scan()

		split := strings.Split(scanner.Text(), " ")

		count, err := strconv.Atoi(split[0])
		must(err)

		P[i] = make([]string, count)

		for j := 0; j < count; j++ {
			P[i][j] = split[1+j]
		}
	}

	err = scanner.Err()
	must(err)

	return Problem{
		M:  M,
		T2: T2,
		T3: T3,
		T4: T4,
		P:  P,
	}
}

func write(filename string, solution Solution) {
	file, err := os.Create(filename)
	defer file.Close()
	must(err)

	_, err = fmt.Fprintln(file, len(solution.D))
	must(err)

	for _, d := range solution.D {
		line := fmt.Sprintf("%d", d.T)

		for _, p := range d.P {
			line += fmt.Sprintf(" %d", p)
		}

		_, err = fmt.Fprintln(file, line)
		must(err)
	}
}

func solve_GetFirst(problem Problem) Solution {
	solution := Solution{
		D: make([]Delivery, 0),
	}

	U := make(map[int]struct{})

	for i := 0; i < problem.T2; i++ {
		D := Delivery{
			T: 2,
			P: make([]int, 0),
		}
		M := make(map[int]struct{})
		for p := 0; p < problem.M; p++ {
			if _, ok := U[p]; ok {
				continue
			}
			if _, ok := M[p]; ok {
				continue
			}

			D.P = append(D.P, p)
			M[p] = struct{}{}

			if len(M) >= 2 {
				break
			}
		}

		if len(D.P) == 2 {
			for k, v := range M {
				U[k] = v
			}
			solution.D = append(solution.D, D)
		}
	}

	for i := 0; i < problem.T3; i++ {
		D := Delivery{
			T: 3,
			P: make([]int, 0),
		}
		M := make(map[int]struct{})
		for p := 0; p < problem.M; p++ {
			if _, ok := U[p]; ok {
				continue
			}
			if _, ok := M[p]; ok {
				continue
			}

			D.P = append(D.P, p)
			M[p] = struct{}{}

			if len(M) >= 3 {
				break
			}
		}

		if len(D.P) == 3 {
			for k, v := range M {
				U[k] = v
			}
			solution.D = append(solution.D, D)
		}
	}

	for i := 0; i < problem.T4; i++ {
		D := Delivery{
			T: 4,
			P: make([]int, 0),
		}
		M := make(map[int]struct{})
		for p := 0; p < problem.M; p++ {
			if _, ok := U[p]; ok {
				continue
			}
			if _, ok := M[p]; ok {
				continue
			}

			D.P = append(D.P, p)
			M[p] = struct{}{}

			if len(M) >= 4 {
				break
			}
		}

		if len(D.P) == 4 {
			for k, v := range M {
				U[k] = v
			}
			solution.D = append(solution.D, D)
		}
	}

	return solution
}

func validate(problem Problem, solution Solution) {
	U := make(map[int]struct{})

	for _, d := range solution.D {
		T := d.T
		P := d.P

		if T != len(P) {
			log.Fatalf("%d person team received %d pizzas", T, len(P))
		}

		for _, p := range P {
			if _, ok := U[p]; ok {
				log.Fatalf("Pizza %d has already been delivered", p)
			}
			U[p] = struct{}{}
		}
	}

	if len(U) > problem.M {
		log.Fatalf("Delivered %d pizzas, %d actually exist", len(U), problem.M)
	}

	log.Printf("Validation success!")
}

func score(problem Problem, solution Solution) {
	score := uint64(0)

	for _, d := range solution.D {
		P := d.P
		I := make(map[string]struct{})

		for _, p := range P {
			for _, i := range problem.P[p] {
				I[i] = struct{}{}
			}
		}

		L := len(I)

		score += uint64(L * L)
	}

	log.Printf("Score %d!", score)
}

func solveA() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\a_example.in")
	solution := solve_GetFirst(problem)

	//log.Printf("%v", solution)

	validate(problem, solution)
	score(problem, solution)
	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\a_example.out", solution)
}

func solveB() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\b_little_bit_of_everything.in")
	solution := solve_GetFirst(problem)

	validate(problem, solution)
	score(problem, solution)

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\b_little_bit_of_everything.out", solution)
}

func solveC() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\c_many_ingredients.in")
	solution := solve_GetFirst(problem)

	validate(problem, solution)
	score(problem, solution)

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\c_many_ingredients.out", solution)
}

func solveD() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\d_many_pizzas.in")
	solution := solve_GetFirst(problem)

	validate(problem, solution)
	score(problem, solution)

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\d_many_pizzas.out", solution)
}

func solveE() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\e_many_teams.in")
	solution := solve_GetFirst(problem)

	validate(problem, solution)
	score(problem, solution)

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\e_many_teams.out", solution)
}

func main() {
	solveA()
	solveB()
	solveC()
	solveD()
	solveE()
}
