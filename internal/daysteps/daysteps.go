package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

/*
	Функция принимает строку с данными, которая содержит количество шагов
	и продолжительность прогулки в формате 3h50m(3 часа 50 минут).

	Возвращает три значения:
		int — количество шагов
		time.Duration — продолжительность прогулки.
		error — ошибку, если что-то пошло не так.

	Функция парсит строку формата "678,0h50m", где:
		678 — количество шагов.
		0h50m — продолжительность прогулки.
*/
func parsePackage(data string) (int, time.Duration, error) {
	dataArr := strings.Split(data, ",")

	if len(dataArr) != 2 {
		return 0, 0, fmt.Errorf("Некоректный формат входного параметра")
	}

	steps, err := strconv.Atoi(dataArr[0])
	if err != nil {
		return 0, 0, err
	}

	duration, err := time.ParseDuration(dataArr[1])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("Кол-во шагов должно быть больше 0")
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("Дистанция должна быть больше 0")
	}

	return steps, duration, nil
}

/*
	У функции три параметра:
		data string — строка с данными, которая содержит количество шагов и продолжительность прогулки в формате 3h50m (3 часа 50 минут).
		weight float64 — вес пользователя в килограммах.
		height float64 — рост пользователя в метрах.
	И одно возвращаемое значение — строка (string).

	Функция должна парсить строку с данными с помощью parsePackage(), вычислять дистанцию в километрах
	и количество потраченных калорий и возвращать строку в таком виде:
		Количество шагов: 792.
		Дистанция составила 0.51 км.
		Вы сожгли 221.33 ккал.
*/
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}

	dist := float64(steps) * stepLength / mInKm

	cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		dist,
		cal,
	)
}
