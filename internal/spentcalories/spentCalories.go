package spentcalories

import (
	"strings"
	"time"
	"strconv"
	"errors"
	"fmt"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// ваш код ниже
	//Разделить строку на слайс строк.
	parts := strings.Split(data, ",")
	//Проверить, чтобы длина слайса была равна 3, так как в строке данных у нас количество шагов, вид активности и продолжительность.
	if len(parts) != 3 {
		return 0, "", 0, errors.New("Неверный формат данных")
	//Преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки. 
	//При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("Неверный формат количества шагов")
	//Преобразовать третий элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration. 
	//Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("Неверный формат продолжительности")
	}
	//Если всё прошло без ошибок, верните количество шагов, вид активности, продолжительность и nil (для ошибки).
	if duration <= 0 {
		return 0, "", 0, errors.New("Продолжительность должна быть больше 0")
	}
	return steps, parts[1], duration, nil
}
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	// ваш код ниже
	return float64(steps) * lenStep / mInKm
}
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	// ваш код ниже
//Проверить, что продолжительность duration больше 0. Если это не так, вернуть 0.
hours := duration.Hours()
	if hours == 0 {
		return 0
	}
//Вычислить дистанцию с помощью distance().
dist := distance(steps)
//Вычислить и вернуть среднюю скорость. Для этого разделите дистанцию на продолжительность в часах. 
//Чтобы перевести продолжительность в часы, воспользуйтесь функцией из пакета time.
return dist / hours
}
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// ваш код ниже
//Получить значения из строки данных с помощью функции parseTraining(), обработать возможные ошибки.
steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return err.Error()
	}
//Проверить, какой вид тренировки был передан в строке, которую парсили (лучше использовать switch). 
//Для каждого из видов тренировки рассчитать дистанцию, среднюю скорость и калории.
dist := distance(steps)
	speed := meanSpeed(steps, duration)

	var calories float64
	switch activity {
	case "running":
		calories = RunningSpentCalories(steps, weight, duration)
	case "walking":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "Неизвестный вид активности"
	}
//Для каждого вида тренировки сформируйте и верните строку, образец которой был представлен выше.
return fmt.Sprintf("Тип тренировки: %s, Длительность: %.2f ч., Дистанция: %.2f км., Скорость: %.2f км/ч, Сожгли калорий: %.2f", activity, duration, dist, speed, calories)
}
//Если был передан неизвестный тип тренировки, верните "неизвестный тип тренировки".

}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// ваш код здесь
//Рассчитать среднюю скорость с помощью meanSpeed().
speed := meanSpeed(steps, duration)
//Рассчитать и вернуть количество калорий. 
//Формула для расчёта: ((runningCaloriesMeanSpeedMultiplier*meanSpeed)-runningCaloriesMeanSpeedShift) * weight
return (runningCaloriesMeanSpeedMultiplier*speed + runningCaloriesMeanSpeedShift) * weight / mInKm * duration.Hours() * minInH
}
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// ваш код здесь
//Рассчитать среднюю скорость с помощью meanSpeed().
speed := meanSpeed(steps, duration)
//Рассчитать и вернуть количество калорий. 
//Формула для расчёта:((walkingCaloriesWeightMultiplier * weight) + (meanSpeed*meanSpeed/height)*walkingSpeedHeightMultiplier) * duration * minInH
return (walkingCaloriesWeightMultiplier*weight + (speed*speed/height)*walkingSpeedHeightMultiplier*weight) * duration.Hours()
}
}
