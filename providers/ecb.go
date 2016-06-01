// Copyright 2016 Openprovider Authors. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package providers

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/openprovider/rates"
	"golang.org/x/text/currency"
)

// ECB represents ECB provider
type ECB struct {
	currencies []currency.Unit
}

const (
	ratesLastURL    = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	rates90daysURL  = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	ratesHistoryURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml"
)

// ECBCurrencies are valid values for currency
var ECBCurrencies = []string{
	AUD, BGN, BRL, CAD, CHF, CNY, CZK, DKK, EUR, GBP, HKD,
	HRK, HUF, IDR, ILS, INR, JPY, KRW, MXN, MYR, NOK,
	NZD, PHP, PLN, RON, RUB, SEK, SGD, THB, TRY, USD, ZAR,
}

// NewECBProvider inits ECB provider record
func NewECBProvider() *ECB {
	ecb := new(ECB)
	// init all units
	for _, unit := range ECBCurrencies {
		if c, err := currency.ParseISO(unit); err == nil {
			ecb.currencies = append(ecb.currencies, c)
		}
	}
	return ecb
}

// ECB XML envelope
type envelope struct {
	Data []struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string `xml:"currency,attr"`
			Rate     string `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

// FetchLast gets exchange rates for the last day
func (ecb *ECB) FetchLast() ([]rates.Rate, []error) {
	return ecb.fetch(ratesLastURL)
}

// Fetch90Days gets exchange rates for 90 days
func (ecb *ECB) Fetch90Days() ([]rates.Rate, []error) {
	return ecb.fetch(rates90daysURL)
}

// FetchHistory gets exchange rates for all existing days
func (ecb *ECB) FetchHistory() ([]rates.Rate, []error) {
	return ecb.fetch(ratesHistoryURL)
}

// FetchLast gets exchange rates for the last day
func (ecb *ECB) fetch(url string) (ecbRates []rates.Rate, errors []error) {
	currentTime := time.Now()
	timeString := currentTime.Format(stdDateTime)

	response, err := http.Get(url)
	if err != nil {
		errors = append(errors, err)
		return
	}
	defer response.Body.Close()

	var raw envelope

	if err := xml.NewDecoder(response.Body).Decode(&raw); err != nil {
		errors = append(errors, err)
		return
	}

	for _, day := range raw.Data {
		if currentTime.Format(stdDate) != day.Date {
			if t, err := time.Parse(stdDate, day.Date); err == nil {
				currentTime = t
				timeString = day.Date + " 00:00:00"
			} else {
				errors = append(errors, err)
			}
		}
		ecbRates = append(ecbRates, rates.Rate{
			Date:           currentTime,
			DateString:     timeString,
			Currency:       currency.EUR,
			CurrencyString: currency.EUR.String(),
			Value:          "1.0000",
		})
		for _, unit := range ecb.currencies {
			for _, item := range day.Rates {
				if item.Currency == unit.String() {
					ecbRates = append(ecbRates, rates.Rate{
						Date:           currentTime,
						DateString:     timeString,
						Currency:       unit,
						CurrencyString: unit.String(),
						Value:          item.Rate,
					})
				}
			}
		}
	}
	return
}
