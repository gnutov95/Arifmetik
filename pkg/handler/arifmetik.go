package handler

import (
	"fmt"
	"sync"
	"time"
)

var (
	Active   int
	NoActive int
	Timers   = Timer{
		TimerPlus:           10,
		TimerMinus:          10,
		TimerMultiply:       10,
		TimerDivide:         10,
		TimerInactiveServer: 1000,
	}
)

// Хранит в себе результат и изначальное выражение
type CalculationResult struct {
	Expression string
	Result     int
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

// Мапа для результатов и выражений
type MapCalculator struct {
	result map[string]Calculator
}
type Timer struct {
	TimerPlus           int
	TimerMinus          int
	TimerMultiply       int
	TimerDivide         int
	TimerInactiveServer int
}

// Функция в который выполняются все вычисления и заполняются все структуры
func Process(expression string, timers Timer) Calculator {
	c := Calculator{}
	a := CalculationResult{Expression: "", Result: 0}

	c.mu.Lock()
	c.countPlus = CountPlus(expression)
	c.countMinus = CountMinus(expression)
	c.countMultiply = CountMultiply(expression)
	c.countDivide = CountDivide(expression)
	c.mu.Unlock()

	go func() {
		Active++
	}()
	result, err := SolveMathExpression(expression)
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
			Active--
			NoActive++
		}()
	case <-time.After(time.Duration(c.Timers.TimerInactiveServer) * time.Second):
		c.Results = CalculationResult{Expression: expression, Result: 2020202020202020}
	}
	return c
}
