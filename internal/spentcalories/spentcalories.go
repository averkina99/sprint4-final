package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

/*
	Функция принимает строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов,
	вид активности и продолжительность активности. Функция возвращает четыре значения:
		int — количество шагов.
		string — вид активности.
		time.Duration — продолжительность активности.
		error — ошибку, если что-то пошло не так.

	Функция парсит строку, переводит данные из строки в соответствующие типы и возвращает эти значения.
*/
func parseTraining(data string) (int, string, time.Duration, error) {
	dataArr := strings.Split(data, ",")

	if len(dataArr) != 3 {
		return 0, "", 0, fmt.Errorf("Некоректный формат входного параметра")
	}

	steps, err := strconv.Atoi(dataArr[0])
	if err != nil {
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(dataArr[2])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Кол-во шагов должно быть больше 0")
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("Дистанция должна быть больше 0")
	}

	return steps, dataArr[1], duration, nil
}

/*
	Функция принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах.
*/
func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}

	return height * stepLengthCoefficient * float64(steps) / mInKm
}

/*
	Функция принимает количество шагов steps, рост пользователя height и продолжительность активности duration
	и возвращает среднюю скорость.
*/
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0.0
	}

	dist := distance(steps, height)

	return dist / duration.Hours()
}

/*
	Функция принимает:
		steps int — количество шагов.
		weight, height float64 — вес(кг.) и рост(м.) пользователя.
		duration time.Duration — продолжительность бега.

	И возвращает два значения:
		float64 — количество калорий, потраченных при беге.
		error — ошибку, если входные параметры некорректны (подумайте, какие значения параметров имеют смысл).
*/
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0.0, fmt.Errorf("Кол-во шагов должно быть больше 0")
	}

	if weight <= 0 {
		return 0.0, fmt.Errorf("Вес должен быть больше 0")
	}

	if height <= 0 {
		return 0.0, fmt.Errorf("Рост должен быть больше 0")
	}

	if duration <= 0 {
		return 0.0, fmt.Errorf("Дистанция должна быть больше 0")
	}

	avgSpeed := meanSpeed(steps, height, duration)

	return (weight * avgSpeed * duration.Minutes()) / minInH, nil
}

/*
	Она совпадает с сигнатурой RunningSpentCalories(). Функция принимает:
		steps int — количество шагов.
		weight, height float64 — вес(кг.) и рост(м.) пользователя.
		duration time.Duration — продолжительность ходьбы.

	И возвращает два значения:
		float64 — количество калорий, потраченных при ходьбе.
		error — ошибку, если входные параметры некорректны (подумайте, какие значения параметров имеют смысл).
*/
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	cal, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0.0, err
	}

	return cal * walkingCaloriesCoefficient, nil
}

/*
	Функция принимает:
		data string — строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов, вид активности и продолжительность активности.
		weight, height float64 — вес (кг.) и рост (м.) пользователя.

	И возвращает два значения:
		string — строка с информацией о тренировке в формате, приведенном ниже.
		error — ошибку, при ее возникновении внутри функции.

	Пример возвращаемой строки:
		Тип тренировки: Бег
		Длительность: 0.75 ч.
		Дистанция: 10.00 км.
		Скорость: 13.34 км/ч
		Сожгли калорий: 18621.75
*/
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, tp, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var cal float64 = 0

	switch tp {
	case "Ходьба":
		cal, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		cal, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		tp,
		duration.Hours(),
		distance(steps, height),
		meanSpeed(steps, height, duration),
		cal,
	), nil
}
