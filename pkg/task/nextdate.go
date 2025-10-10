package task

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func daysNextDate(now time.Time, dstartTime time.Time, params string) (time.Time, error) {
	days, err := strconv.Atoi(string(params))
	if err != nil {
		return now, err
	}
	if days <= 0 || days > 400 {
		return now, fmt.Errorf("invalid days count: %d", days)
	}
	// Сразу добавляем заданный интервал для следующей даты
	dstartTime = dstartTime.AddDate(0, 0, days)
	// Проверка на текущую дату, если про
	for dstartTime.Before(now) || dstartTime.Equal(now) {
		dstartTime = dstartTime.AddDate(0, 0, days)
	}

	return dstartTime, nil
}

func weeksNextDate(now time.Time, dstartTime time.Time, params string) (time.Time, error) {
	if len(params) == 0 {
		return now, fmt.Errorf("unexpected weekdays count: %d", len(params))
	}

	weekdays := strings.Split(params, ",")
	// Перегоняем строку в слайс интов для удобства
	intWeekdays := []int{}
	for _, str := range weekdays {
		num, err := strconv.Atoi(str)
		if err != nil {
			return now, err
		}
		if num < 1 || num > 7 {
			return now, fmt.Errorf("weekday %d out of range", num)
		}
		intWeekdays = append(intWeekdays, num)
	}
	// Порядок дней в params может быть случайным?
	slices.Sort(intWeekdays)

	startTime := dstartTime
	// Если заданная дата уже наступила - ищем ближайший день недели
	if dstartTime.Before(now) {
		startTime = now
	}
	currentWeekday := int(startTime.Weekday())

	for _, weekday := range intWeekdays {
		if currentWeekday < weekday {
			return startTime.AddDate(0, 0, weekday-currentWeekday), nil
		}
	}
	return startTime.AddDate(0, 0, 7-currentWeekday+intWeekdays[0]), nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var nextTime time.Time
	dstartTime, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	if len(repeat) == 0 {
		return "", nil
	}

	rule := strings.Split(repeat, " ")

	switch rule[0] {
	case "d":
		if len(rule) != 2 {
			return "", fmt.Errorf("unexpected repeate rule len: %d", len(rule))
		}
		nextTime, err = daysNextDate(now, dstartTime, rule[1])
		if err != nil {
			return "", err
		}

	case "y":
		dstartTime = dstartTime.AddDate(1, 0, 0)

		for dstartTime.Before(now) {
			dstartTime = dstartTime.AddDate(1, 0, 0)
		}
		nextTime = dstartTime

	case "w":
		nextTime, err = weeksNextDate(now, dstartTime, rule[1])
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("repeate rule %s is not allowed", rule[0])
	}

	return nextTime.Format("20060102"), nil
}
