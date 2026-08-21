// Package envconfig предоставляет вспомогательные функции для загрузки
// конфигурации приложения из переменных окружения.
package envconfig

import (
	"fmt"
	"os"
	"strings"
)

// Required возвращает очищенное от пробелов значение переменной окружения key.
// Если переменная не установлена или содержит только пробелы, возвращается ошибка.
func Required(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)

	if !ok || value == "" {
		return "", fmt.Errorf("environment variable %s is required", key)
	}

	return value, nil
}

// OrDefault возвращает очищенное от пробелов значение переменной окружения key.
// Если переменная не установлена или содержит только пробелы,
// возвращается defaultValue.
func OrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)

	if !ok || value == "" {
		return defaultValue
	}

	return value
}
