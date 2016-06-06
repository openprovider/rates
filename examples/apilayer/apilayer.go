package main

import (
	"fmt"

	"github.com/openprovider/rates"
	"github.com/openprovider/rates/providers"
)

func main() {
	// Get exchange rates for today
	service := rates.New(
		// any collection of providers which implement rates.Provider interface
		providers.NewAPILayerProvider(
			&rates.Options{
				Token: "xxx",
				Currencies: []string{
					providers.EUR,
					providers.USD,
					providers.CHF,
					providers.HKD,
				},
			},
		),
	)
	items, errors := service.FetchLast()
	if len(errors) != 0 {
		fmt.Println(errors)
	}
	fmt.Println(service.Name(), "exchange rates for today")
	for index, item := range items {
		fmt.Printf("%d. %s - %v\r\n", index+1, item.Currency, item.Value)
	}

	// Get historical exchange rates
	settings := make(map[string]interface{})
	settings["date"] = "2015-01-01"
	registry := rates.Registry{
		providers.NewAPILayerProvider(
			&rates.Options{
				Token: "xxx",
				Currencies: []string{
					providers.GHS,
					providers.LKR,
				},
				Settings: settings,
			},
		),
	}
	data, errors := registry.FetchHistory()
	if len(errors) != 0 {
		fmt.Println(errors)
	}
	fmt.Println(registry.Name(), "exchange rates for", settings["date"])
	for index, item := range data {
		fmt.Printf("%d. %s - %v\r\n", index+1, item.Currency, item.Value)
	}
}
