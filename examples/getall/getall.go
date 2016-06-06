package main

import (
	"fmt"

	"github.com/openprovider/rates"
	"github.com/openprovider/rates/providers"
)

func main() {
	service := rates.New(
		// any collection of providers which implement rates.Provider interface
		providers.NewECBProvider(new(rates.Options)),
	)
	rates, errors := service.FetchLast()
	if len(errors) != 0 {
		fmt.Println(errors)
	}
	fmt.Println(service.Name(), "exchange rates for today")
	for index, rate := range rates {
		fmt.Printf("%d. %s - %v\r\n", index+1, rate.Currency, rate.Value)
	}
}
