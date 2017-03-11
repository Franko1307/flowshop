package main

import "fmt"

func RemoveFromSlice(s []int, i int) []int {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

type Element struct {
	times        []int
	waiting_time int
	tardiness    int
}

type Calendar struct {
	elements []Element
	fitness  int
}

type Machine struct {
	elements []*Element
	time     int
	ind_time int
	occupied bool
}

type Line struct {
	machines []Machine
	line     int
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (line Line) add(element *Element) {

	best := line.machines[0].time
	idx := 0
	for i := range line.machines {
		if line.machines[i].time < best {
			best = line.machines[i].time
			idx = i
		}
	}

	if !line.machines[idx].occupied {
		line.machines[idx].occupied = true
		line.machines[idx].ind_time = element.times[line.line]
	}

	line.machines[idx].elements = append(line.machines[idx].elements, element)
	line.machines[idx].time += element.times[line.line]

}

func pass_time(lines []*Line) {

	for i := range lines {
		for j := range lines[i].machines {
			if lines[i].machines[j].ind_time > 0 {
				lines[i].machines[j].ind_time -= 1
			}
			if lines[i].machines[j].time > 0 {
				lines[i].machines[j].time -= 1
			}
		}
	}

}

func (line Line) remove_from_line() []*Element {

	elements_to_remove := make([]*Element, 0)
	for i := range line.machines {
		if line.machines[i].occupied && line.machines[i].ind_time == 0 {
			elements_to_remove = append(elements_to_remove, line.machines[i].elements[0])
			line.machines[i].elements[0] = nil //<-- optional ?
			line.machines[i].elements = line.machines[i].elements[1:]

			if line.machines[i].time > 0 {
				line.machines[i].ind_time = line.machines[i].elements[0].times[line.line]
			} else {
				line.machines[i].ind_time = 0
				line.machines[i].occupied = false
			}

		}
	}

	return elements_to_remove

}
func elements_in_lines(lines []*Line) bool {
	for i := range lines {
		for j := range lines[i].machines {
			if len(lines[i].machines[j].elements) > 0 {
				return true
			}
		}
	}
	return false
}

func ComputeFitness(calendar Calendar, machines []int) (int, int) {

	lines := make([]*Line, 0)

	for i := range machines {
		lines = append(lines, &Line{machines: make([]Machine, machines[i]), line: i})
	}

	for i := range calendar.elements {
		lines[0].add(&calendar.elements[i])
	}

	fitness := 0
	tardiness := 0
	aux_elements := make([]*Element, 0)

	for elements_in_lines(lines) {
		fitness++

		pass_time(lines)

		for i := range lines {
			aux_elements = append(aux_elements, lines[i].remove_from_line()...)
			for j := 0; j < len(aux_elements); j++ {
				if i+1 == len(lines) {
					for x := range aux_elements {
						tardiness += Max(fitness-aux_elements[x].tardiness, 0)
					}
					aux_elements = aux_elements[:0]
				} else {
					lines[i+1].add(aux_elements[i])
					aux_elements[i] = nil
					aux_elements = aux_elements[1:]
				}
			}
		}
	}

	return fitness, tardiness

}

func greedy(idx int, elements []Element, machines []int) {

	fmt.Println("Añadiendo tarea ", idx)
	calendar := Calendar{elements: make([]Element, 0)}
	calendar.elements = append(calendar.elements, Element{times: elements[idx].times, tardiness: elements[idx].tardiness})

	fitness, _ := ComputeFitness(calendar, machines)
	fmt.Println("Fitness: ", fitness)
	rest := make([]int, 0)
	for i := range elements {
		rest = append(rest, i)
	}

	rest = RemoveFromSlice(rest, idx)
	var aux int
	for len(rest) != 1 {
		bestFitness, _ := ComputeFitness(Calendar{elements: append(calendar.elements, Element{times: elements[rest[0]].times, tardiness: elements[rest[0]].tardiness})}, machines)
		auxidx := rest[0]
		fmt.Println("Con tarea ", rest[0], " tendriamos fitness de ", bestFitness)
		for i := 1; i < len(rest); i++ {
			aux, _ = ComputeFitness(Calendar{elements: append(calendar.elements, Element{times: elements[rest[i]].times, tardiness: elements[rest[i]].tardiness})}, machines)
			fmt.Println("Con tarea ", rest[i], " tendriamos fitness de ", aux)
			if aux < bestFitness {
				auxidx = rest[i]
				bestFitness = aux
			}
		}
		fmt.Println("Añadiendo tarea ", auxidx, " porque tendremos fitness de ", bestFitness)
		calendar.elements = append(calendar.elements, Element{times: elements[auxidx].times, tardiness: elements[auxidx].tardiness})
		rest = RemoveFromSlice(rest, idx)
	}

	fmt.Println("Añadiendo última tarea: ", rest[0])
	calendar.elements = append(calendar.elements, Element{times: elements[rest[0]].times, tardiness: elements[rest[0]].tardiness})
	fitness, _ = ComputeFitness(calendar, machines)
	fmt.Println("Fitness final: ", fitness)
}

func main() {

	e := make([]Element, 3)

	e[0].times = []int{2, 1}
	e[0].tardiness = 3
	e[1].times = []int{3, 4}
	e[1].tardiness = 8
	e[2].times = []int{3, 2}
	e[2].tardiness = 11
	greedy(1, e, []int{1, 2})

}
