package main

import (
	"fmt"
	"sync"
)

// Counter - структура счётчика с методами для работы с ним
type Counter struct {
	value     int           // текущее значение счётчика
	target    int           // целевое значение счётчика
	numGoroutines int       // количество горутин для обновления счётчика
	mutex     sync.Mutex    // мьютекс для синхронизации доступа к счётчику
	wg        sync.WaitGroup // группа ожидания для горутин
	done      chan bool     // канал для сигнализации завершения
}

// NewCounter - конструктор для создания нового счётчика
func NewCounter(target int, numGoroutines int) *Counter {
	return &Counter{
		value:         0,
		target:        target,
		numGoroutines: numGoroutines,
		done:          make(chan bool),
	}
}

// Increment - метод для увеличения значения счётчика на 1
func (c *Counter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if c.value < c.target { // проверяем, что счётчик ещё не достиг цели
		c.value++
		fmt.Printf("Текущее значение счётчика: %d\n", c.value)
		
		// проверяем, достиг ли счётчик целевого значения
		if c.value >= c.target {
			close(c.done) // закрываем канал, чтобы сигнализировать о завершении
		}
	}
}

// GetValue - метод для получения текущего значения счётчика
func (c *Counter) GetValue() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.value
}

// GetTarget - метод для получения целевого значения счётчика
func (c *Counter) GetTarget() int {
	return c.target
}

// IsReached - метод для проверки, достигнута ли целевая величина
func (c *Counter) IsReached() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.value >= c.target
}

// Start - метод для запуска горутин, которые будут увеличивать счётчик
func (c *Counter) Start() {
	// запускаем указанное количество горутин
	for i := 0; i < c.numGoroutines; i++ {
		c.wg.Add(1)
		go func(goroutineID int) {
			defer c.wg.Done()
			
			// каждая горутина будет увеличивать счётчик, пока тот не достигнет цели
			for {
				select {
				case <-c.done: // если достигнуто целевое значение, выходим из цикла
					fmt.Printf("Горутина %d завершена, цель достигнута\n", goroutineID)
					return
				default:
					// проверяем, достигли ли мы цели перед инкрементом
					if c.IsReached() {
						fmt.Printf("Горутина %d завершена, цель уже достигнута\n", goroutineID)
						return
					}
					
					// увеличиваем счётчик
					c.Increment()
					
					// небольшая задержка, чтобы другие горутины могли работать
					// в реальном приложении можно использовать time.Sleep(), но без импорта time
					// просто продолжаем выполнение
				}
			}
		}(i)
	}
	
	// ждём завершения всех горутин
	c.wg.Wait()
}

// Reset - метод для сброса счётчика к начальному значению
func (c *Counter) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value = 0
	// пересоздаём канал done, так как он был закрыт
	c.done = make(chan bool)
}

func main() {
	// пример использования
	fmt.Println("Программа счётчика с использованием горутин")
	fmt.Println("===========================================")
	
	// создаем счётчик с целевым значением 10 и 3 горутинами
	counter := NewCounter(10, 3)
	
	fmt.Printf("Запуск счётчика с целевым значением: %d, количество горутин: %d\n", 
		counter.GetTarget(), counter.numGoroutines)
	
	// запускаем процесс подсчёта
	counter.Start()
	
	fmt.Printf("\nФинальное значение счётчика: %d\n", counter.GetValue())
	fmt.Println("Процесс завершён успешно!")
	
	// демонстрация динамической проверки достижения цели
	fmt.Printf("Цель достигнута: %t\n", counter.IsReached())
}