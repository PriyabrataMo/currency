package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/huh"
	_ "github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func isNum(input string) error {
	if _, err := strconv.ParseFloat(input, 64); err != nil {
		return errors.New("invalid number")
	}
	return nil
}

type Currency struct {
	Rates struct {
		USD float64 `json:"USD"`
		GBP float64 `json:"GBP"`
		INR float64 `json:"INR"`
		JPY float64 `json:"JPY"`
		EUR float64 `json:"EUR"`
	} `json:"Rates"`
}

var (
	cur1    string = "INR"
	cur2    string = "INR"
	amt     string
	confirm bool
)

func main() {

	// Create a new form

	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))
	dynamic1 := "INR"

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Currency Converter").
			Description("Welcome to CurrencyConv™.\n\nYour GoTo place to convert Currency\n\n").
			Next(true).
			NextLabel("Next"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose the currency you want to convert from").
				Options(
					huh.NewOptions("INR", "USD", "EUR", "GBP", "JPY")...).
				Value(&cur1), // store the chosen option in the "burger" variable

			huh.NewSelect[string]().
				Value(&cur2).
				Title("Choose the currency you want to convert to").
				TitleFunc(func() string {
					if cur1 == "" {
						return "Choose the currency you want to convert to" // Default title if cur1 is empty
					}
					dynamic1 = cur1
					return fmt.Sprintf("Choose the currency you want %s convert to", cur1)
				}, &cur1).
				//Title(fmt.Sprintf("Choose the currency you want to convert to")).
				Options(
					huh.NewOptions("INR", "USD", "EUR", "GBP", "JPY")...),
		),
		huh.NewGroup(
			huh.NewInput().
				TitleFunc(func() string {

					return fmt.Sprintf("Type in the amount you want to convert %s to %s", dynamic1, cur2)
				}, &cur2).
				Prompt("?").
				Validate(isNum).
				Value(&amt),

			huh.NewConfirm().
				Title("Are you sure?").
				Negative("Wait, No.").
				Affirmative("Yes!").
				Value(&confirm).
				WithTheme(huh.ThemeBase()),
		),
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
	if confirm {
		fmt.Println("Please make sure you're picking the right option, okay? 🙄")
		return
	}
	result := make(chan float64)
	done := make(chan bool)
	go fetchAnswer(result, done)

	select {
	case <-done:
		break
	default:
		loading(1)
	}
	finalAnswer := <-result
	out := fmt.Sprintf("\n%s %s is %0.2f %s", amt, cur1, finalAnswer, cur2)
	fmt.Println(out)

}

func loading(s int) {

	chars := spinner.Globe.Frames
	timer := time.NewTicker(time.Second * time.Duration(s)) // Stop after 's' seconds.
	defer timer.Stop()

	// Start the spinner loop
	for {
		select {
		case <-timer.C: // After 's' seconds, stop the spinner.
			fmt.Println("\r")
			return
		default:
			for _, c := range chars {
				fmt.Printf("\r%s", string(c))     // Overwrite previous spinner frame
				time.Sleep((time.Second * 1) / 5) // Adjust spinner speed
			}
		}
	}

}
func fetchAnswer(result chan float64, done chan bool) {
	err := godotenv.Load()
	apikey := os.Getenv("API_KEY")
	if apikey == "" {
		log.Fatal("API_KEY is not set")
	}
	res, err := http.Get("https://api.exchangeratesapi.io/v1/latest?access_key=" + apikey)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if res.StatusCode != 200 {
		panic("Currency Api call failed with status code: " + strconv.Itoa(res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var currency Currency

	err = json.Unmarshal(body, &currency)
	if err != nil {
		panic(err)
	}
	Rates := currency.Rates
	ratemap := make(map[string]float64)
	ratemap["USD"], ratemap["GBP"], ratemap["INR"], ratemap["JPY"], ratemap["EUR"] = Rates.USD, Rates.GBP, Rates.INR, Rates.JPY, Rates.EUR

	rate := ratemap[cur2] / ratemap[cur1]
	f, _ := strconv.ParseFloat(amt, 64)
	var out = rate * f
	result <- out
	done <- true

}
