package main

import (
	"fmt"
	"os"
	"webcrawler/db"
	"webcrawler/search"
)

func main() {
	var query string

	// Запрос у пользователя
	fmt.Println("Введите запрос для поиска:")
	fmt.Scanln(&query)

	// Проверка на пустой запрос
	if query == "" {
		fmt.Println("Ошибка: запрос не может быть пустым.")
		os.Exit(1)
	}

	// Соединение с базой данных
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Парсим ссылки по запросу
	links, err := search.SearchLinks(query)
	if err != nil {
		fmt.Println("Ошибка при поиске:", err)
		os.Exit(1)
	}

	// Сохраняем ссылки в базе данных
	for _, link := range links {
		err = db.SaveLink(conn, link)
		if err != nil {
			fmt.Println("Ошибка при сохранении ссылки:", err)
		} else {
			fmt.Println("Ссылка сохранена:", link)
		}
	}
}
