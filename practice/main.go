package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

type Problem struct {
	M  int        // Number of pizzas.
	T2 int        // Number of 2 person teams.
	T3 int        // Number of 3 person teams.
	T4 int        // Number of 4 person teams.
	P  [][]string // Pizza ingredients.
}

type Delivery struct {
	T int   // Number of people in the team.
	P []int // Pizza indexes for that team.
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

func solve_MostIngredients(problem Problem) Solution {
	// Assigns pizzas with the most number of ingredients to larger teams first.
	solution := Solution{
		D: make([]Delivery, 0),
	}
	// P - pizza ID, C - number of ingredients on that pizza.
	C := make([]struct {
		P int
		C int
	}, len(problem.P))
	for idx, p := range problem.P {
		C[idx].P = idx
		C[idx].C = len(p)
	}
	// Sort C on pizza ingredient count in descending order.
	sort.Slice(C, func(i, j int) bool {
		c1 := C[i].C
		c2 := C[j].C
		return c1 > c2
	})

	// Total number of pizzas delivered.
	H := 0
	// Iterate over the 4 person team and deliver pizzas with many of ingredients first.
	for idx := 0; idx < problem.T4; idx++ {
		D := Delivery{
			T: 4,
			P: make([]int, 0),
		}
		// Check for when we have fewer pizzas left than the team requires.
		Hmax := min(H+4, len(C))
		for H < Hmax {
			D.P = append(D.P, C[H].P)
			H++
		}
		// Check that each team member will receive a pizza or the whole team doesn't get any pizzas.
		if len(D.P) == 4 {
			solution.D = append(solution.D, D)
		}
	}
	// Same but for 3 person team.
	for idx := 0; idx < problem.T3; idx++ {
		D := Delivery{
			T: 3,
			P: make([]int, 0),
		}

		Hmax := min(H+3, len(C))
		for H < Hmax {
			D.P = append(D.P, C[H].P)
			H++
		}

		if len(D.P) == 3 {
			solution.D = append(solution.D, D)
		}
	}
	// Same but for 2 person team.
	for idx := 0; idx < problem.T2; idx++ {
		D := Delivery{
			T: 2,
			P: make([]int, 0),
		}

		Hmax := min(H+2, len(C))
		for H < Hmax {
			D.P = append(D.P, C[H].P)
			H++
		}

		if len(D.P) == 2 {
			solution.D = append(solution.D, D)
		}
	}

	return solution
}

func initialSolution(problem Problem) Solution {
	return solve_MostIngredients(problem)
}

func solve_Anneal(problem Problem) Solution {
	solution := initialSolution(problem)

	Tmax := 100.0
	Tmin := 0.5
	// Current temperature.
	Ti := Tmax

	for Ti > Tmin {
		// Score before random swaps.
		S1 := score(problem, solution)
		// Total number of deliveries.
		Dl := len(solution.D)
		// Pick 2 random deliveries indexes to swap.
		Didx1 := rand.Intn(Dl)
		Didx2 := rand.Intn(Dl)
		// Numbers of pizzas in each delivery.
		Pl1 := len(solution.D[Didx1].P)
		Pl2 := len(solution.D[Didx2].P)
		// Random pizza indexes from each delivery to swap.
		Pidx1 := rand.Intn(Pl1)
		Pidx2 := rand.Intn(Pl2)
		// Swap 2 pizzas.
		P1 := solution.D[Didx1].P[Pidx1]
		P2 := solution.D[Didx2].P[Pidx2]

		solution.D[Didx1].P[Pidx1] = P2
		solution.D[Didx2].P[Pidx2] = P1
		// Get the resulting score.
		S2 := score(problem, solution)

		Sdelta := S1 - S2

		if Sdelta <= 0 {
			// The original score was smaller than what we got after the swap, so keep the current state.
		} else {
			// Got a lower score than we started with, so decide if we still want to keep it.
			P := math.Exp(float64(-Sdelta) / Ti)
			R := rand.Float64()

			if R <= P {
				// Still want to swap.
			} else {
				// Put back original pizzas.
				solution.D[Didx1].P[Pidx1] = P1
				solution.D[Didx2].P[Pidx2] = P2
			}
		}
		// Decrease the temperature.
		Ti = Ti - (Ti * 0.001)

		log.Printf("S=%d Ti=%f", score(problem, solution), Ti)
	}

	return solution
}

func validate(problem Problem, solution Solution) {
	// Map of all delivered pizzas.
	D := make(map[int]struct{})

	for _, d := range solution.D {
		T := d.T
		P := d.P

		if T != len(P) {
			log.Fatalf("%d person team received %d pizzas (each team member must receive a pizza)", T, len(P))
		}

		for _, p := range P {
			if _, ok := D[p]; ok {
				log.Fatalf("Pizza %d has already been delivered", p)
			}
			D[p] = struct{}{}
		}
	}

	if len(D) > problem.M {
		log.Fatalf("Delivered more pizzas than actually exist %d vs %d", len(D), problem.M)
	}

	log.Printf("Validation success!")
}

func score(problem Problem, solution Solution) int64 {
	S := int64(0)

	for _, d := range solution.D {
		// Map of unique ingredients for a given delivery.
		I := make(map[string]struct{})

		for _, p := range d.P {
			for _, i := range problem.P[p] {
				I[i] = struct{}{}
			}
		}
		// Score is the total number of unique ingredients squared.
		L := len(I)
		S += int64(L * L)
	}

	return S
}

func solveA() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\a_example.in")
	solution := solve_Anneal(problem)

	validate(problem, solution)
	log.Printf("Score %d!", score(problem, solution))

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\a_example.out", solution)
}

func solveB() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\b_little_bit_of_everything.in")
	solution := solve_Anneal(problem)

	validate(problem, solution)
	log.Printf("Score %d!", score(problem, solution))

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\b_little_bit_of_everything.out", solution)
}

func solveC() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\c_many_ingredients.in")
	solution := solve_Anneal(problem)

	validate(problem, solution)
	log.Printf("Score %d!", score(problem, solution))

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\c_many_ingredients.out", solution)
}

func solveD() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\d_many_pizzas.in")
	solution := solve_Anneal(problem)

	validate(problem, solution)
	log.Printf("Score %d!", score(problem, solution))

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\d_many_pizzas.out", solution)
}

func solveE() {
	problem := read("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\input\\e_many_teams.in")
	solution := solve_Anneal(problem)

	validate(problem, solution)
	log.Printf("Score %d!", score(problem, solution))

	write("C:\\Users\\Juris Liepins\\go\\src\\github.com\\jurisliepins\\google-hash-code-2021\\practice\\output\\e_many_teams.out", solution)
}

func main() {
	//solveA()
	//solveB()
	//solveC()
	//solveD()
	//solveE()
}
