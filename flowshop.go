package main

import "fmt"

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
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

func (line Line) add_to_line(element *Element) {

	best := line.machines[0].time
	idx := 0
	for i := 0; i < len(line.machines); i++ {
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

	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i].machines); j++ {
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
	for i := 0; i < len(line.machines); i++ {
		if line.machines[i].occupied && line.machines[i].ind_time == 0 {
			fmt.Println("Apunto de remover de la linea ", 0)
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
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i].machines); j++ {
			if len(lines[i].machines[j].elements) > 0 {
				return true
			}
		}
	}

	return false
}

func printlines(lines []*Line) {
	for i := 0; i < len(lines); i++ {
		fmt.Println("Linea: ", i, lines[i])
		for j := 0; j < len(lines[i].machines); j++ {
			fmt.Println("Maquina ", j, lines[i].machines[j], "tiempo: ", lines[i].machines[j].time, "ind time : ", lines[i].machines[j].ind_time)
			for k := 0; k < len(lines[i].machines[j].elements); k++ {
				fmt.Println("Element: ", k, lines[i].machines[j].elements[k])
			}
		}
	}
}

func (calendar Calendar) get_fitness(n []int) (int, int) {

	lines := make([]*Line, len(n))
	for i := 0; i < len(n); i++ {
		lines[i] = &Line{}
		lines[i].line = i
		lines[i].machines = make([]Machine, n[i])
	}
	fmt.Println("--------")
	printlines(lines)
	fmt.Println("--------")
	for i := 0; i < len(calendar.elements); i++ {
		lines[0].add_to_line(&calendar.elements[i])
	}
	fmt.Println("&--------&")
	printlines(lines)
	fmt.Println("&--------&")

	aux_elements := make([]*Element, 0)
	fitness := 0
	tardiness := 0
	for elements_in_lines(lines) {

		fitness += 1
		fmt.Println("Fitness: ", fitness)
		if fitness == 12 {
			return 0, 0
		}
		pass_time(lines)
		fmt.Println("El tiempo pasa...")
		for i := 0; i < len(lines); i++ {
			fmt.Println("Antes: ")
			printlines(lines)
			aux_elements = append(aux_elements, lines[i].remove_from_line()...)
			fmt.Println("longitud: ", len(aux_elements))
			fmt.Println("Despues de sacar pero antes de meter")
			printlines(lines)
			println("")
			println("")
			for j := 0; j < len(aux_elements); j++ {
				if i+1 == len(lines) {
					fmt.Println("Saliendo elementos...")
					for x := 0; x < len(aux_elements); x++ {
						fmt.Println(aux_elements[0])
						tardiness += Max(fitness-aux_elements[x].tardiness, 0)
					}
					fmt.Println("tardiness: ", tardiness)
					fmt.Println("Done")
					aux_elements = aux_elements[:0]
				} else {
					lines[i+1].add_to_line(aux_elements[i])
					aux_elements[i] = nil
					aux_elements = aux_elements[1:]
				}
			}
			fmt.Println("@@@@@@@@@@@@")
			fmt.Println("Despues de sacar y despues de meter")
			printlines(lines)
			fmt.Println("@@@@@@")

		}
	}

	return fitness, tardiness
}

func randomCalendar(n int, m int) Calendar {

	calendar := Calendar{elements: make([]Element, n)}

	for i := 0; i < n; i++ {
		calendar.elements[i].times = make([]int, m)
	}

	calendar.elements[0].times[0] = 3
	calendar.elements[0].times[1] = 2
	calendar.elements[0].tardiness = 11
	calendar.elements[1].times[0] = 1
	calendar.elements[1].times[1] = 2
	calendar.elements[1].tardiness = 3
	calendar.elements[2].times[0] = 3
	calendar.elements[2].times[1] = 4
	calendar.elements[2].tardiness = 8

	return calendar
}

func main() {

	calendar := randomCalendar(3, 2)

	fmt.Println(calendar)

	fmt.Println(calendar.get_fitness([]int{1, 2}))

	fmt.Println(calendar)
}
