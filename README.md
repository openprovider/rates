Exchange rates provider
=======================

A package used to manage exchange rates from any provider

[![GoDoc](https://godoc.org/github.com/openprovider/rates?status.svg)](https://godoc.org/github.com/openprovider/rates)

### Examples

Get all exchange rates for the ECB provider
```go
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
```

Get exchange rates for EUR, USD, CHF, HKD
```go
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
```

## Authors

[Igor Dolzhikov](https://github.com/takama)

## Contributors

All the contributors are welcome. If you would like to be the contributor please accept some rules.
- The pull requests will be accepted only in "develop" branch
- All modifications or additions should be tested

Thank you for your understanding!

## License

[MIT Public License](https://github.com/openprovider/rates/blob/master/LICENSE)
