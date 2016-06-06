// Copyright 2016 Openprovider Authors. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

/*
Package rates 0.2.1
This package helps to manage exchange rates from any provider

Example 1: Get all exchange rates for the ECB Provider

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
        fmt.Println("European Central Bank exchange rates for today")
        for index, rate := range rates {
            fmt.Printf("%d. %s - %v\r\n", index+1, rate.Currency, rate.Value)
        }
    }

Example 2: Get exchange rates for EUR, USD, CHF, HKD

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
        fmt.Println("European Central Bank exchange rates for today")
        for index, rate := range rates {
            fmt.Printf("%d. %s - %v\r\n", index+1, rate.Currency, rate.Value)
        }
    }

Exchange Rates Provider
*/
package rates

import (
	"strings"
	"time"

	"golang.org/x/text/currency"
)

// Rate represent date and currency exchange rates
type Rate struct {
	ID       uint64        `json:"id,omitempty"`
	Date     string        `json:"date"`
	Currency string        `json:"currency"`
	Time     time.Time     `json:"-"`
	Base     currency.Unit `json:"-"`
	Unit     currency.Unit `json:"-"`
	Value    interface{}   `json:"value"`
}

// Options is some specific things for the specific provider
// It should configure the provider to manage currencies
type Options struct {
	// API key/token
	Token string
	// List of the currencies which need to get from the provider
	// If it is empty, should get all of existing currencies from the provider
	Currencies []string
	// Flexible settings list
	Settings map[string]interface{}
}

// Provider holds methods for providers which implement this interface
type Provider interface {
	Name() string
	FetchLast() (rates []Rate, errors []error)
	FetchHistory() (rates []Rate, errors []error)
}

// Registry contains registered providers
type Registry []Provider

// New service which contains registered providers
func New(providers ...Provider) Provider {
	var registry Registry
	for _, provider := range providers {
		registry = append(registry, provider)
	}
	return registry
}

// Name returns name of the provider
func (registry Registry) Name() string {
	var names []string
	for _, provider := range registry {
		names = append(names, provider.Name())
	}
	return "Registry: " + strings.Join(names, ", ")
}

// FetchLast returns exchange rates from all registered providers on last day
func (registry Registry) FetchLast() (rates []Rate, errors []error) {
	for _, provider := range registry {
		r, errs := provider.FetchLast()
		rates = append(rates, r...)
		errors = append(errors, errs...)
	}
	return
}

// FetchHistory returns exchange rates from all registered providers from history
func (registry Registry) FetchHistory() (rates []Rate, errors []error) {
	for _, provider := range registry {
		r, errs := provider.FetchHistory()
		rates = append(rates, r...)
		errors = append(errors, errs...)
	}
	return
}
