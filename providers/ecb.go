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
	ecbName       = "European Central Bank"
	ecbLastURL    = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	ecb90daysURL  = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	ecbHistoryURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml"
)

// ECBCurrencies are valid types of currencies for that provider
var ECBCurrencies = []string{
	AUD, BGN, BRL, CAD, CHF, CNY, CZK, DKK, EUR, GBP, HKD,
	HRK, HUF, IDR, ILS, INR, JPY, KRW, MXN, MYR, NOK,
	NZD, PHP, PLN, RON, RUB, SEK, SGD, THB, TRY, USD, ZAR,
}

// NewECBProvider inits ECB provider record
func NewECBProvider(options *rates.Options) *ECB {
	ecb := new(ECB)
	// init all units
	if len(options.Currencies) == 0 {
		options.Currencies = append(options.Currencies, ECBCurrencies...)
	}
	for _, unit := range options.Currencies {
		if c, err := currency.ParseISO(unit); err == nil {
			ecb.currencies = append(ecb.currencies, c)
		}
	}
	return ecb
}

// ECB XML envelope
type ecbEnvelope struct {
	Data []struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string `xml:"currency,attr"`
			Rate     string `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

// Name returns name of the provider
func (ecb *ECB) Name() string {
	return ecbName
}

// FetchLast gets exchange rates for the last day
func (ecb *ECB) FetchLast() ([]rates.Rate, []error) {
	return ecb.fetch(ecbLastURL)
}

// Fetch90Days gets exchange rates for 90 days
func (ecb *ECB) Fetch90Days() ([]rates.Rate, []error) {
	return ecb.fetch(ecb90daysURL)
}

// FetchHistory gets exchange rates for all existing days
func (ecb *ECB) FetchHistory() ([]rates.Rate, []error) {
	return ecb.fetch(ecbHistoryURL)
}

// FetchLast gets exchange rates for the last day
func (ecb *ECB) fetch(url string) (ecbRates []rates.Rate, ecbErrors []error) {
	currentTime := time.Now()
	date := currentTime.Format(stdDateTime)

	response, err := http.Get(url)
	if err != nil {
		ecbErrors = append(ecbErrors, err)
		return
	}
	defer response.Body.Close()

	var raw ecbEnvelope

	if err := xml.NewDecoder(response.Body).Decode(&raw); err != nil {
		ecbErrors = append(ecbErrors, err)
		return
	}

	for _, day := range raw.Data {
		if currentTime.Format(stdDate) != day.Date {
			if t, err := time.Parse(stdDate, day.Date); err == nil {
				currentTime = t
				date = t.Format(stdDateTime)
			} else {
				ecbErrors = append(ecbErrors, err)
			}
		}
		for _, unit := range ecb.currencies {
			if unit == currency.EUR {
				ecbRates = append(ecbRates, rates.Rate{
					Time:     currentTime,
					Date:     date,
					Base:     currency.EUR,
					Unit:     currency.EUR,
					Currency: currency.EUR.String() + "/" + currency.EUR.String(),
					Value:    "1.0000",
				})
			}
			for _, item := range day.Rates {
				if item.Currency == unit.String() {
					ecbRates = append(ecbRates, rates.Rate{
						Time:     currentTime,
						Date:     date,
						Base:     currency.EUR,
						Unit:     unit,
						Currency: currency.EUR.String() + "/" + unit.String(),
						Value:    item.Rate,
					})
				}
			}
		}
	}
	return
}
