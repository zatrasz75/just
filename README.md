# Just Logger

Простой и удобный логгер для Go с поддержкой цветного вывода.

## Установка

```bash
go get github.com/zatrasz75/just
````

## Использование
###  Базовое логирование
```bash
package main

import (
  "fmt"
  "github.com/zatrasz75/just/logger"
)
  
const logFilePath = "./var/log/main.log"

func main() {
   // Если logFilePath оставить "" то записи логов в файл не будет, только в консоль
	l, err := logger.NewLogger(logFilePath)
		if err != nil {
		fmt.Println("Ошибка при создании файла логгера:", err)
		return
	}
	defer l.Close()
	
	l.Info("Это информационное сообщение", "key1", "value1")
	l.Error("Это информационное сообщение", fmt.Errorf("что-то пошло не так"))
}
```

### Уровни логирования
- Info (зеленый)

- Success (синий)

- Trace (серый)

- Error (красный)

- Warn (желтый)

- Fatal (красный)

- Debug (голубой)

- Critical (фиолетовый)

- Panic (оранжевый)

- Security (яркий фиолетовый)

### Логирование HTTP-запросов (middleware)
```bash
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zatrasz75/just/logger"
)

func main() {
	// Создаем логгер для вывода в консоль
	consoleLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	respLogger := logger.NewResponLogger(consoleLogger)

	// Создаем маршрутизатор
	router := mux.NewRouter()

	// Используем middleware для логирования запросов
	router.Use(respLogger.LoggingResponse)

	// Пример маршрута
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Just Logger!"))
	})

	// Запуск сервера
	http.ListenAndServe(":8080", router)
}

```
### Пример вывода логов в терминале

```bash
0000/00/00 00:00:00.000000 >-- ip: 127.0.0.1:62066, host: localhost:8080 url: /users/login, method: POST status code: 200 OK, trace id: e0543b67-eaa3-4367-861e-24600bae2924
```