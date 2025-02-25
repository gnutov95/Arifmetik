package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Data struct {
	Id     int
	Input  string
	Result int
	t      bool
}

var (
	arrCalculator = MapCalculator{result: map[string]Calculator{}}
	count         = 0
	arr1          = make(map[int]Data)
	keys          = make([]int, 0) // Убрал емкость из make, так как она будет динамически изменяться
	mutex         = sync.Mutex{}
)

func UpdateFunc(w http.ResponseWriter, r *http.Request) {
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
}

func TimeUpdate(w http.ResponseWriter, r *http.Request) {
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
}

func MainPage(w http.ResponseWriter, r *http.Request) {
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
}
func Switch2(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("htmlDirectory/page3.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct {
		Act   int
		NoAct int
	}{
		Act:   Active,
		NoAct: NoActive,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
