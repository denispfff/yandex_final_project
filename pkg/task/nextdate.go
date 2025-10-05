package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

func weeksNextDate(now time.Time, params string) (time.Time, error) {
	if len(params) == 0 {
		return now, nil
	}

	weekdays := strings.Split(params, ",")
	nowWeekday := int(now.Weekday())

	intWeekdays := make([]int, len(weekdays))
	for _, str := range weekdays {
		num, err := strconv.Atoi(str)
		if err != nil {
			return now, err
		}
		intWeekdays = append(intWeekdays, num)
	}

	for _, weekday := range intWeekdays {
		if nowWeekday < weekday {
			return now.AddDate(0, 0, weekday-nowWeekday), nil
		}
	}
	return now.AddDate(0, 0, 7+intWeekdays[0]-nowWeekday), nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var nextTime time.Time
	dstartTime, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}
	fmt.Println(repeat)
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
		nextTime, err = weeksNextDate(now, rule[1])
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("repeate rule %s is not allowed", rule[0])
	}

	return nextTime.Format("20060102"), nil
}
