package Timekeeper

import (
	"chicha/pkg/config"
	"chicha/pkg/data"
	"fmt"
)

// Обрабатывает данные гонки в памяти и возвращает обновленные круги и возможную ошибку.
func calculateRaceInMemory(currentTimekeeperTask Data.RawData, previousLaps []Data.Lap, config Config.Settings) (updatedLaps []Data.Lap, err error) {
	var currentLap Data.Lap
	var lapIsValid bool

	// Проверка валидности текущих данных круга.
	lapIsValid, err = checkLapIsValid(currentTimekeeperTask, previousLaps, config)
	if err != nil {
		fmt.Printf("Ошибка валидации данных: %s \n", err)
		return
	} else if !lapIsValid {
		fmt.Println("Получены невалидные данные круга. Пропускаем.")
		updatedLaps = previousLaps
		return
	}

	// Инициализация текущего круга
	currentLap, err = initializeCurrentLap(currentTimekeeperTask, previousLaps, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Обновление текущего круга на основе предыдущих данных
	updatedLaps, err = updateCurrentLap(currentLap, previousLaps, currentTimekeeperTask, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Рассчитываем позиции и отставания
	calculatePositionsAndDelays(updatedLaps, config)

	return updatedLaps, nil
}


