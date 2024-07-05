package shared

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
)

// Experiment configurations
const WarmupMessages int = 100000

// Controllers
const BasicPID string = "BasicPID"
const DeadzonePID string = "DeadzonePID"
const HPA string = "HPA"
const ASTAR string = "ASTAR"
const ErrorSquareProportional = "ErrorSquareProportional"
const ErrorSquareFull = "ErrorSquareFull"
const IncrementalPID = "IncrementalPID"
const SetPointWeighting = "SetPointWeighting"
const SmoothingPID = "SmoothingPID"
const FuzzyController = "FuzzyController"

// Fuzzy Logic
type OutputX struct {
	Mx  []float64
	Out []float64
}

// Utils functions
func ErrorHandler(f string, msg string) {
	fmt.Println(f + "::" + msg)
	os.Exit(0)
}
func GetFunction() string {
	fpcs := make([]uintptr, 1)

	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		fmt.Println("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		fmt.Println("MSG CALLER WAS NIL")
	}

	// Print the file name and line number
	//fmt.Println(caller.FileLine(fpcs[0]-1))

	// Print the name of the function
	//fmt.Println(caller.Name())

	return caller.Name()
}
func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}
func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func ReadConfParameter(pName string) string {
	params := make(map[string]string)

	// open & load file
	filePath := "C:\\Users\\user\\go\\adaptive-moms\\controllers"
	fileName := "config.yaml"

	readFile, err := os.Open(filePath + "\\" + fileName)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	// Load existing parameters
	for _, line := range fileLines {
		p := strings.Split(line, "=")
		params[p[0]] = p[1]
	}

	// Check if pName exist in the list of parameters
	pValue, ok := params[pName]
	if !ok {
		ErrorHandler(GetFunction(), "Controller parameter does not exist.")
	}

	return pValue
}
