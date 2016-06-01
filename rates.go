// Copyright 2016 Openprovider Authors. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

/*
Package rates 0.0.1
This package helps to manage exchange rates from any provider

Example 1:

    package main

    import (
        "fmt"

        "github.com/openprovider/rates"
        "github.com/openprovider/rates/providers"
    )

    func main() {
        registry := rates.Registry{
            // any collection of providers which implement rates.Provider interface
            providers.NewECBProvider(),
        }
        service := rates.New(registry)
        rates, errors := service.FetchLast()
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
	"time"

	"golang.org/x/text/currency"
)

// Rate represent date and currency exchange rates
type Rate struct {
	ID             uint64        `json:"id,omitempty"`
	DateString     string        `json:"date"`
	Date           time.Time     `json:"-"`
	CurrencyString string        `json:"currency"`
	Currency       currency.Unit `json:"-"`
	Value          interface{}   `json:"value"`
}

// Provider holds methods for providers which implement this interface
type Provider interface {
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
