package main

import (
	"fmt"

	"github.com/openprovider/rates"
	"github.com/openprovider/rates/providers"
)

func main() {
	registry := rates.Registry{
		// any collection of providers which implement rates.Provider interface
		providers.NewECBProvider(
			&rates.Options{
				Currencies: []string{
					providers.EUR,
					providers.USD,
					providers.CHF,
					providers.HKD,
				},
			},
		),
	}
	rates, errors := registry.FetchLast()
	if len(errors) != 0 {
		fmt.Println(errors)
	}
	fmt.Println(registry.Name(), "exchange rates for today")
	for index, rate := range rates {
		fmt.Printf("%d. %s - %v\r\n", index+1, rate.Currency, rate.Value)
	}
}
