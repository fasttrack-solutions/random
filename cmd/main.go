package main

import (
	"bufio"
	"fmt"
	"github.com/fasttrack-solutions/random"
	"github.com/nexidian/gocliselect"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	menu := gocliselect.NewMenu("Select test to run")

	menu.AddItem("Random Uniform Float64 with range (0,1]", "UniformFloat64")
	menu.AddItem("Random Uniform Int with range (min, max)", "UniformInt64")
	menu.AddItem("Linear to Exponential transformation", "TransformToExponential")
	menu.AddItem("Truncate", "Truncate")
	menu.AddItem("Round", "Round")
	menu.AddItem("Exit", "Exit")

	choice := menu.Display()

	err := os.MkdirAll("results", 0775)
	if err != nil {
		fmt.Println("error creating results directory")
		return
	}

	switch choice {
	case "UniformFloat64":
		for {
			res, errPrompt := stringPrompt("how many results do you want to generate?")
			if errPrompt != nil {
				fmt.Println("error getting results from prompt:", errPrompt)
			} else {
				numbersToGenerate, errAtoi := strconv.Atoi(res)
				if errAtoi != nil {
					fmt.Println(res, "is an invalid number, try again")
				} else {
					fmt.Print("generating ", res, " random numbers... ")

					fileName, errGenerate := generateUniformFloat64(numbersToGenerate)
					if errGenerate != nil {
						fmt.Print("error: ", errGenerate)
						return
					}

					fmt.Println("file with results was generated:", fileName)

					break
				}
			}
		}

	case "UniformInt64":
		for {
			res, errPrompt := stringPrompt("how many results do you want to generate?")
			if errPrompt != nil {
				fmt.Println("error getting results from prompt:", errPrompt)
			} else {
				num, errAtoi := strconv.Atoi(res)
				if errAtoi != nil {
					fmt.Println(res, "is an invalid number, try again")
				} else {
					res, errPrompt = stringPrompt("whats the minimum number?")
					if errPrompt != nil {
						fmt.Println("error getting results from prompt:", errPrompt)
					} else {
						minimumNumber, errAtoiMin := strconv.Atoi(res)
						if errAtoiMin != nil {
							fmt.Println(res, "is an invalid number, try again")
						} else if minimumNumber < 0 {
							fmt.Println(res, "is less than 0, try again")
						} else {
							res, errPrompt = stringPrompt("whats the maximum number?")
							if errPrompt != nil {
								fmt.Println("error getting results from prompt:", errPrompt)
							} else {
								maximumNumber, errAtoiMax := strconv.Atoi(res)
								if errAtoiMax != nil {
									fmt.Println(res, "is an invalid number, try again")
								} else if maximumNumber <= minimumNumber {
									fmt.Println("maximum cannot be less than or equal to minimum number, try again")
								} else {

									fmt.Print("generating ", res, " random numbers between ", minimumNumber, " and ", maximumNumber, "... ")

									fileName, errGenerate := generateUniformInt64(num, int64(minimumNumber), int64(maximumNumber))
									if errGenerate != nil {
										fmt.Print("error: ", errGenerate)
										return
									}

									fmt.Println("file with results was generated:", fileName)

									break
								}
							}
						}
					}
				}
			}
		}

	case "Exit":
		return

	default:
		fmt.Println("oops... selected test is not yet implemented")
	}

	// generateRandom(100)
}

func stringPrompt(label string) (string, error) {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		_, err := fmt.Fprint(os.Stderr, label+" ")
		if err != nil {
			return "", err
		}
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s), nil
}

func generateUniformFloat64(numbersToGenerate int) (string, error) {
	fileName := fmt.Sprintf("results/UniformFloat64-%v.csv", time.Now().UnixMilli())

	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer f.Close()

	for i := 0; i < numbersToGenerate; i++ {
		rnd, errRnr := random.UniformFloat64()
		if errRnr != nil {
			return "", errRnr
		}

		_, errWriteString := f.WriteString(fmt.Sprintf("%v\n", rnd))
		if errWriteString != nil {
			return "", errWriteString
		}
	}

	return fileName, nil
}

func generateUniformInt64(numbersToGenerate int, min int64, max int64) (string, error) {
	fileName := fmt.Sprintf("results/UniformInt64-%v.csv", time.Now().UnixMilli())

	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer f.Close()

	for i := 0; i < numbersToGenerate; i++ {
		rnd, errRnr := random.UniformInt64(min, max)
		if errRnr != nil {
			return "", errRnr
		}

		_, errWriteString := f.WriteString(fmt.Sprintf("%v\n", rnd))
		if errWriteString != nil {
			return "", errWriteString
		}
	}

	return fileName, nil
}
