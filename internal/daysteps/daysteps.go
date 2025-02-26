package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alex1234ak/fitness-4-sprint/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, " ")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid data format: expected steps and duration")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("invalid steps format")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, errors.New("invalid duration format")
	}

	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	if steps <= 0 {
		fmt.Println("Error: steps must be greater than 0")
		return ""
	}

	distanceMeters := float64(steps) * StepLength
	distanceKilometers := distanceMeters / spentcalories.MetersInKilometer

	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	result := fmt.Sprintf("Steps: %d, Distance: %.2f km, Time: %s, Calories: %.2f kcal",
		steps, distanceKilometers, duration, calories)

	return result
}
