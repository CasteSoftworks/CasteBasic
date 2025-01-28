package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Instruction struct {
	command string
	args    []string
}

var program = map[int]Instruction{}
var variables = map[string]int{}
var stringVar = map[string]string{}

func loadProgramFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Errore nell'apertura del file %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parseLine(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Errore durante la lettura del file: %v\n", err)
		os.Exit(1)
	}
}

func parseLine(line string) {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '\'' { // Salta linee vuote o commenti
		return
	}

	parts := strings.Fields(line)
	lineNumber, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("Errore: il numero di linea non è valido: %s\n", parts[0])
		return
	}

	command := strings.ToUpper(parts[1])
	args := parts[2:]
	program[lineNumber] = Instruction{command: command, args: args}

	//fmt.Println("\t", lineNumber, program[lineNumber])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Utilizzo: go run main.go <nomefile.bas>")
		os.Exit(1)
	}

	filename := os.Args[1]
	loadProgramFromFile(filename)
	executeProgram()
}

func executeProgram() {
	currentLine := getFirstLine()

	for currentLine > 0 {
		instr := program[currentLine]
		switch instr.command {
		case "PRINT":
			handlePrint(instr.args)
		case "LET":
			handleLet(instr.args)
		case "GOTO":
			line, _ := strconv.Atoi(instr.args[0])
			currentLine = line
			continue
		case "IF":
			handleIf(instr.args, &currentLine)
			continue
		case "INPUT":
			handleInput(instr.args)
		case "SINPUT":
			handleStringInput(instr.args)
		case "END":
			fmt.Println("Fine del programma.")
			return
		default:
			fmt.Printf("Errore: comando sconosciuto %s\n", instr.command)
		}
		currentLine = getNextLine(currentLine)
	}
}

func handleStringInput(s []string) {
	if len(s) != 1 {
		fmt.Println("Errore di sintassi in SINPUT")
		os.Exit(1)
	}
	if s[0][1] != '$' {
		fmt.Println("Errore: variabile non di tipo stringa")
		os.Exit(1)
	}
	varName := s[0]
	fmt.Print("? ")
	var value string
	fmt.Scan(&value)
	stringVar[varName] = value
}

func handleInput(s []string) {
	if len(s) != 1 {
		fmt.Println("Errore di sintassi in INPUT")
		os.Exit(1)
	}
	varName := s[0]
	fmt.Print("? ")
	var value int
	_, err := fmt.Scanf("%d", &value)
	if err != nil {
		fmt.Println("Errore: inserire un numero intero")
		os.Exit(1)
	}
	variables[varName] = value
}

/*func handlePrint(args []string) {
	inQuotes := false
	var buffer string

	for _, arg := range args {
		if inQuotes {
			if arg[len(arg)-1] == '"' {
				buffer += " " + arg[:len(arg)-1]
				fmt.Println(buffer)
				inQuotes = false
				buffer = ""
			} else {
				buffer += " " + arg
			}
		} else {
			if len(arg) > 2 && arg[0] == '"' {
				if arg[len(arg)-1] == '"' {
					fmt.Println(arg[1 : len(arg)-1])
				} else {
					inQuotes = true
					buffer = arg[1:]
				}
			} else {
				result := evaluateExpression([]string{arg})
				fmt.Print(result)
			}
		}
	}

	if inQuotes {
		fmt.Println(buffer)

	}
}*/

func handlePrint(args []string) {
	inQuotes := false
	var buffer string

	for _, arg := range args {
		if inQuotes {
			if arg[len(arg)-1] == '"' {
				buffer += " " + arg[:len(arg)-1]
				fmt.Println(buffer)
				inQuotes = false
				buffer = ""
			} else {
				buffer += " " + arg
			}
		} else {
			if len(arg) > 2 && arg[0] == '"' {
				if arg[len(arg)-1] == '"' {
					fmt.Println(arg[1 : len(arg)-1])
				} else {
					inQuotes = true
					buffer = arg[1:]
				}
			} else if len(arg) > 1 && arg[len(arg)-1] == '$' {
				// Stampa una variabile stringa (x$)
				varName := arg
				if val, exists := stringVar[varName]; exists {
					fmt.Print(val)
				} else {
					fmt.Printf("Errore: variabile stringa %s non inizializzata\n", varName)
				}
			} else {
				// Stampa il risultato di un'espressione numerica
				result := evaluateExpression([]string{arg})
				fmt.Println(result)
			}
		}
	}

	if inQuotes {
		fmt.Println(buffer)
	}
}

func handleLet(args []string) {
	if len(args) < 3 || args[1] != "=" {
		fmt.Println("Errore di sintassi in LET")
		os.Exit(1)
	}
	varName := args[0]
	value := evaluateExpression(args[2:])
	variables[varName] = value
}

func handleIf(args []string, currentLine *int) {
	thenIndex := indexOf(args, "THEN")
	elseIndex := indexOf(args, "ELSE")

	if thenIndex == -1 {
		fmt.Println("Errore di sintassi in IF: manca THEN")
		return
	}

	condition := args[:thenIndex]
	if evaluateCondition(condition) {
		executeInlineCommand(args[thenIndex+1:], currentLine)
	} else if elseIndex != -1 {
		executeInlineCommand(args[elseIndex+1:], currentLine)
	}
}

func executeInlineCommand(args []string, currentLine *int) {
	if len(args) == 0 {
		fmt.Println("Errore: comando mancante")
		return
	}
	switch args[0] {
	case "PRINT":
		handlePrint(args[1:])
	case "LET":
		handleLet(args[1:])
	case "GOTO":
		line, _ := strconv.Atoi(args[1])
		*currentLine = line
	default:
		fmt.Printf("Errore: comando sconosciuto %s\n", args[0])
	}
}

func evaluateCondition(condition []string) bool {
	if len(condition) != 3 {
		fmt.Println("Errore di sintassi nella condizione")
		return false
	}
	left := evaluateExpression([]string{condition[0]})
	right := evaluateExpression([]string{condition[2]})
	switch condition[1] {
	case "=":
		return left == right
	case "<":
		return left < right
	case ">":
		return left > right
	case "<=":
		return left <= right
	case ">=":
		return left >= right
	case "<>":
		return left != right
	}
	fmt.Println("Operatore condizionale sconosciuto:", condition[1])
	return false
}

func evaluateExpression(expr []string) int {
	if len(expr) == 1 {
		// Se è un numero, restituiscilo
		if val, err := strconv.Atoi(expr[0]); err == nil {
			return val
		}
		// Se è una variabile, verifica se è definita
		if val, exists := variables[expr[0]]; exists {
			return val
		}
		// Variabile non definita
		fmt.Printf("Errore: variabile %s non inizializzata\n", expr[0])
		os.Exit(1)
	}
	if len(expr) == 3 {
		left := evaluateExpression([]string{expr[0]})
		right := evaluateExpression([]string{expr[2]})
		switch expr[1] {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			if right == 0 {
				fmt.Println("Errore: divisione per zero")
				os.Exit(1)
			}
			return left / right
		case "MOD":
			if right == 0 {
				fmt.Println("Errore: divisione per zero")
				os.Exit(1)
			}
			return left % right
		}
	}
	fmt.Println("Errore nell'espressione:", strings.Join(expr, " "))
	os.Exit(1)
	return 0 // Non raggiunto
}

func indexOf(arr []string, word string) int {
	for i, v := range arr {
		if v == word {
			return i
		}
	}
	return -1
}

func getFirstLine() int {
	lines := make([]int, 0, len(program))
	for line := range program {
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return 0
	}
	sort.Ints(lines) // Ordina le linee
	return lines[0]
}

func getNextLine(currentLine int) int {
	lines := make([]int, 0, len(program))
	for line := range program {
		lines = append(lines, line)
	}
	sort.Ints(lines)
	for _, line := range lines {
		if line > currentLine {
			return line
		}
	}
	return 0 // Fine del programma
}
