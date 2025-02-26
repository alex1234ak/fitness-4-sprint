package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                            = 0.65  // средняя длина шага.
	mInKm                              = 1000  // количество метров в километре.
	minInH                             = 60    // количество минут в часе.
	kmhInMsec                          = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM                              = 100   // количество сантиметров в метре.
	MetersInKilometer                  = 1000
	runningCaloriesMeanSpeedMultiplier = 18.0  // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0  // среднее количество сжигаемых калорий при беге.
	walkingCaloriesWeightMultiplier    = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier       = 0.029 // множитель роста.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid data format: expected steps, activity, and duration")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("invalid steps format")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("invalid duration format")
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be greater than 0")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

func meanSpeed(steps int, duration time.Duration) float64 {
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return distance(steps) / hours
}

func TrainingInfo(data string, weight, height float64) string {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return err.Error()
	}

	dist := distance(steps)
	speed := meanSpeed(steps, duration)

	var calories float64
	switch activity {
	case "running":
		calories = RunningSpentCalories(steps, weight, duration)
	case "walking":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "unknown activity type"
	}

	return fmt.Sprintf("Training type: %s, Duration: %.2f h, Distance: %.2f km, Speed: %.2f km/h, Calories burned: %.2f",
		activity, duration.Hours(), dist, speed, calories)
}

func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	return (runningCaloriesMeanSpeedMultiplier*speed - runningCaloriesMeanSpeedShift) * weight / mInKm * duration.Hours() * minInH
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	if height == 0 {
		return 0
	}
	speed := meanSpeed(steps, duration)
	return (walkingCaloriesWeightMultiplier*weight + (speed*speed/height)*walkingSpeedHeightMultiplier) * duration.Hours()
}
