package main

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
)

func main() {
	// URL Selenium Server
	seleniumURL := "http://localhost:4444/wd/hub"

	// Настройки браузера
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// Подключение к Selenium
	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		fmt.Printf("Ошибка подключения к Selenium: %v\n", err)
		return
	}
	defer wd.Quit()

	// Переход на страницу IMDb
	if err := wd.Get("https://www.imdb.com/search/title/?title_type=feature"); err != nil {
		fmt.Printf("Ошибка открытия страницы IMDb: %v\n", err)
		return
	}

	// Принятие куки
	acceptButton, err := wd.FindElements(selenium.ByCSSSelector, "button[data-testid='accept-button']")
	if err == nil && len(acceptButton) > 0 {
		err := acceptButton[0].Click()
		if err != nil {
			fmt.Printf("Ошибка при клике на кнопку 'Accept': %v\n", err)
			return
		}
		fmt.Println("Принято использование куки.")
		time.Sleep(1 * time.Second)
	}

	// Парсинг фильмов
	for {
		// Извлекаем названия фильмов
		filmElements, err := wd.FindElements(selenium.ByCSSSelector, "a.ipc-title-link-wrapper h3.ipc-title__text")
		if err != nil {
			fmt.Printf("Ошибка при извлечении названий фильмов: %v\n", err)
			return
		}
		fmt.Println("Текущие фильмы:")
		for _, element := range filmElements {
			title, err := element.Text()
			if err != nil {
				fmt.Printf("Ошибка при извлечении текста: %v\n", err)
				return
			}
			fmt.Println(title)
		}

		// Проверяем наличие кнопки "50 more"
		buttons, err := wd.FindElements(selenium.ByCSSSelector, "button.ipc-see-more__button")
		if err != nil {
			fmt.Printf("Ошибка поиска кнопки '50 more': %v\n", err)
			return
		}
		if len(buttons) > 0 {
			_, err := wd.ExecuteScript(`arguments[0].scrollIntoView({behavior: "smooth", block: "center"});`, []interface{}{buttons[0]})
			if err != nil {
				fmt.Printf("Ошибка при прокрутке страницы: %v\n", err)
				return
			}
			time.Sleep(2 * time.Second)

			_, err = wd.ExecuteScript(`arguments[0].click();`, []interface{}{buttons[0]})
			if err != nil {
				fmt.Printf("Ошибка при клике на кнопку '50 more': %v\n", err)
				return
			}
			fmt.Println("Нажата кнопка '50 more', загружаются новые фильмы...")
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("Кнопка '50 more' не найдена. Завершаем парсинг.")
			break
		}
	}
}
