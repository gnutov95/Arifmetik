package main

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Хранит в себе результат и изначальное выражение
type CalculationResult struct {
	Expression string
	Result     int
}

// Мапа для результатов и выражений
type MapCalculator struct {
	result map[string]Calculator
}

// Хранит в себе все  данные о выражении и его выполнении
type Calculator struct {
	Results       CalculationResult // Результаты
	countPlus     int               // кол-во плюсов
	countMinus    int               // кол-во минусов
	countMultiply int               // кол-во знаков умножения
	countDivide   int               // кол-во знаков деления

	tComlete bool
	Timers   Timer // установленное время выполнения
	mu       sync.Mutex
}

type Timer struct {
	TimerPlus           int
	TimerMinus          int
	TimerMultiply       int
	TimerDivide         int
	TimerInactiveServer int
}

var active int
var noActive int

// методы проверок на наличие символов
func hasPlus(input string) bool {
	return strings.Contains(input, "+")
}

func hasMinus(input string) bool {
	return strings.Contains(input, "-")
}

func hasMultiply(input string) bool {
	return strings.Contains(input, "*")
}

func hasDivide(input string) bool {
	return strings.Contains(input, "/")
}

// запись имеющихся символов
func hasSymbol(input string) ([]int, bool) {
	f := false
	i := make([]int, 4)
	if hasPlus(input) {
		if i[0] == 0 {
			i[0] = 1
			f = true
		}
	}
	if hasMinus(input) {
		if i[1] == 0 {
			i[1] = 2
			f = true
		}
	}
	if hasMultiply(input) {
		if i[2] == 0 {
			i[2] = 3
			f = true
		}
	}
	if hasDivide(input) {
		if i[3] == 0 {
			i[3] = 4
			f = true
		}
	}
	return i, f
}

// метод вычисления выражения
func solveMathExpression(expression string) (float64, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	finalResult, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("unable to convert result to float64")
	}

	return finalResult, nil
}

// хранит в себе стандартное время выполенений операций
var Timers = Timer{
	TimerPlus:           10,
	TimerMinus:          10,
	TimerMultiply:       10,
	TimerDivide:         10,
	TimerInactiveServer: 1000,
}

// Функция для подсчета символа '+'
func countPlus(input string) int {
	count := 0
	for _, char := range input {
		if char == '+' {
			count++
		}
	}
	return count
}

// Функция для подсчета символа '-'
func countMinus(input string) int {
	count := 0
	for _, char := range input {
		if char == '-' {
			count++
		}
	}
	return count
}

// Функция для подсчета символа '*'
func countMultiply(input string) int {
	count := 0
	for _, char := range input {
		if char == '*' {
			count++
		}
	}
	return count
}

// Функция для подсчета символа '/'
func countDivide(input string) int {
	count := 0
	for _, char := range input {
		if char == '/' {
			count++
		}
	}
	return count
}

// Функция в который выполняются все вычисления и заполняются все структуры
func Process(expression string, timers Timer) Calculator {
	c := Calculator{}
	a := CalculationResult{Expression: "", Result: 0}

	c.mu.Lock()
	c.countPlus = countPlus(expression)
	c.countMinus = countMinus(expression)
	c.countMultiply = countMultiply(expression)
	c.countDivide = countDivide(expression)
	c.mu.Unlock()

	go func() {
		active++
	}()
	result, err := solveMathExpression(expression)
	if err != nil {
		fmt.Println("Error processing expression:", err)
		return Calculator{}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.Timers = timers
	a = CalculationResult{Expression: expression, Result: int(result)}
	c.Results = a
	c.tComlete = false
	fmt.Println(c.Timers.TimerMinus*c.countMinus + c.Timers.TimerPlus*c.countPlus + c.Timers.TimerMultiply*c.countMultiply + c.countDivide*c.countDivide)

	select { // Задержка с использованием таймера
	case <-time.After(time.Duration(c.Timers.TimerMinus*c.countMinus+c.Timers.TimerPlus*c.countPlus+c.Timers.TimerMultiply*c.countMultiply+c.countDivide*c.countDivide) * time.Second):
		go func() {
			active--
			noActive++
		}()
	case <-time.After(time.Duration(c.Timers.TimerInactiveServer) * time.Second):
		c.Results = CalculationResult{Expression: expression, Result: 2020202020202020}
	}
	return c
}

type Data struct {
	Id     int
	Input  string
	Result int
	t      bool
}

func main() {

	arrCalculator := MapCalculator{result: map[string]Calculator{}}
	count := 0
	arr1 := make(map[int]Data)
	keys := make([]int, 0) // Убрал емкость из make, так как она будет динамически изменяться
	mutex := sync.Mutex{}  // Добавил мьютекс для безопасного доступа к данным
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		// Подготавливаем данные для отправки
		data := struct {
			Keys []int
			Data map[int]Data
		}{
			Keys: keys,
			Data: arr1,
		}

		// Отправляем данные в формате JSON
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/update_table_body2_data", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		// Подготавливаем данные для отправки
		data := struct {
			Keys []int
			Data map[int]Data
		}{
			Keys: keys,
			Data: arr1,
		}

		// Отправляем данные в формате JSON
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			inputData := r.Form.Get("inputData")
			if inputData != "" {
				fmt.Println("Received data:", inputData)

				// Добавляем данные в первую таблицу сразу
				mutex.Lock()
				arr1[count] = Data{
					Id:     count + 1,
					Input:  inputData,
					Result: 0, // Значение результата по умолчанию
					t:      false,
				}
				keys = append(keys, count)
				count++
				mutex.Unlock()

				go func(inputData string) {
					// Обработка данных и обновление результатов
					c := Process(inputData, Timers) // Передача таймеров в функцию Process
					mutex.Lock()
					defer mutex.Unlock()
					arrCalculator.result[inputData] = c
					// Обновляем результат в первой таблице
					for i, d := range arr1 {
						if d.Input == inputData {
							arr1[i] = Data{
								Id:     d.Id,
								Input:  d.Input,
								Result: c.Results.Result,
								t:      c.tComlete,
							}
						}
					}
				}(inputData) // Передача inputData в качестве аргумента в анонимную функцию
			}
		}
		tmpl, err := template.ParseFiles("htmlDirectory/page1.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		err = tmpl.Execute(w, struct {
			Keys []int
			Data map[int]Data
		}{
			Keys: keys,
			Data: arr1,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
		var mu sync.Mutex
		if r.Method == "POST" {
			r.ParseForm()
			input1, err1 := strconv.Atoi(r.Form.Get("input1"))
			input2, err2 := strconv.Atoi(r.Form.Get("input2"))
			input3, err3 := strconv.Atoi(r.Form.Get("input3"))
			input4, err4 := strconv.Atoi(r.Form.Get("input4"))
			input5, err5 := strconv.Atoi(r.Form.Get("input5"))
			if err1 != nil && err2 != nil && err3 != nil && err4 != nil && err5 != nil {
				log.Fatalln("Problems Strings to int")
			}
			if input1 != 0 && input2 != 0 && input3 != 0 && input4 != 0 && input5 != 0 {

				mu.Lock()
				Timers = Timer{
					TimerPlus:           input1,
					TimerMinus:          input2,
					TimerMultiply:       input3,
					TimerDivide:         input4,
					TimerInactiveServer: input5,
				}
				mu.Unlock()
			}
		}
		fmt.Println(Timers.TimerMinus)

		tmpl, err := template.ParseFiles("htmlDirectory/page2.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/switch2", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("htmlDirectory/page3.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, struct {
			Act   int
			NoAct int
		}{
			Act:   active,
			NoAct: noActive,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
