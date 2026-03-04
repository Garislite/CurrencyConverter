package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func getRate(currency string) float64 {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://www.cbr.ru/scripts/XML_daily.asp", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	text := string(body)

	valutes := strings.Split(text, "<Valute")

	for _, v := range valutes {
		if !strings.Contains(v, currency) {
			continue
		}

		s := strings.Index(v, "<Value>")
		if s == -1 {
			continue
		}
		e := strings.Index(v[s:], "</Value>")
		if e == -1 {
			continue
		}

		valStr := v[s+7 : s+e]
		valStr = strings.Replace(valStr, ",", ".", -1)
		rate, _ := strconv.ParseFloat(valStr, 64)

		ns := strings.Index(v, "<Nominal>")
		if ns != -1 {
			ne := strings.Index(v[ns:], "</Nominal>")
			if ne != -1 {
				nomStr := v[ns+9 : ns+ne]
				nom, _ := strconv.ParseFloat(nomStr, 64)
				if nom > 0 {
					return rate / nom
				}
			}
		}
		return rate
	}
	return 0
}

func main() {
	currencies := []string{"USD", "EUR", "CNY", "GBP", "TRY"}
	rates := make(map[string]float64)

	for _, code := range currencies {
		rate := getRate(code)
		if rate > 0 {
			rates[code] = rate
			fmt.Printf("1 %s = %.2f RUB\n", code, rate)
		}
	}

	var rub float64
	var toCur string
	var again string

	for {
		fmt.Print("\nEnter RUB amount: ")
		fmt.Scan(&rub)

		fmt.Print("Convert to (USD, EUR, CNY, GBP, TRY): ")
		fmt.Scan(&toCur)
		toCur = strings.ToUpper(toCur)

		if rate, ok := rates[toCur]; ok && rate > 0 {
			fmt.Printf("= %.2f %s\n", rub/rate, toCur)
		} else {
			fmt.Println("Invalid currency")
		}

		fmt.Print("Again? (y/n): ")
		fmt.Scan(&again)
		if again == "n" {
			break
		}
	}
}
