package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const apiUrl = "http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s"

type WeatherData struct {
	Main struct {
		Temperature float64 `json:"temp"`
		Humidity    int     `json:"humidity"`
	} `json:"main"`
	Name string `json:"name"`
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	loadEnv()
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")

	vars := mux.Vars(r)
	city := vars["city"]
	url := fmt.Sprintf(apiUrl, city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Ошибка при выполнении запроса", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Ошибка во время запроса к API OpenWeatherMap. Код состояния: "+resp.Status, http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении данных из ответа", http.StatusInternalServerError)
		return
	}

	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Ошибка при разборе данных JSON", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Погода в городе %s:\n", data.Name)
	fmt.Fprintf(w, "Температура: %.2f°C\n", data.Main.Temperature-273.15)
	fmt.Fprintf(w, "Влажность: %d%%\n", data.Main.Humidity)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello)))")
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/weather/{city}", weatherHandler)
	r.HandleFunc("/weather", helloHandler)

	http.Handle("/", r)

	log.Println("Сервер запущен на :8080...")
	http.ListenAndServe(":8080", nil)
}
