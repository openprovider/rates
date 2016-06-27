// Copyright 2016 Openprovider Authors. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package providers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/openprovider/rates"
	"golang.org/x/text/currency"
)

// APILayer represents multi currency provider
type APILayer struct {
	currencies []currency.Unit
	date       string
	lastURL    string
	historyURL string
}

const (
	apiLayerName       = "API Layer"
	apiLayerLastURL    = "http://apilayer.net/api/live"
	apiLayerHistoryURL = "http://apilayer.net/api/historical"
)

// APILayerCurrencies are valid types of currencies for that provider
var APILayerCurrencies = []string{
	AED, AFN, ALL, AMD, ANG, AOA, ARS, AUD, AWG, AZN, BAM, BBD, BDT, BGN,
	BHD, BIF, BMD, BND, BOB, BRL, BSD, BTC, BTN, BWP, BYR, BZD, CAD, CDF,
	CHF, CLF, CLP, CNY, COP, CRC, CUP, CVE, CZK, DJF, DKK, DOP, DZD, EEK,
	EGP, ERN, ETB, EUR, FJD, FKP, GBP, GEL, GGP, GHS, GIP, GMD, GNF, GTQ,
	GYD, HKD, HNL, HRK, HTG, HUF, IDR, ILS, IMP, INR, IQD, IRR, ISK, JEP,
	JMD, JOD, JPY, KES, KGS, KHR, KMF, KPW, KRW, KWD, KYD, KZT, LAK, LBP,
	LKR, LRD, LSL, LTL, LVL, LYD, MAD, MDL, MGA, MKD, MMK, MNT, MOP, MRO,
	MUR, MVR, MWK, MXN, MYR, MZN, NAD, NGN, NIO, NOK, NPR, NZD, OMR, PAB,
	PEN, PGK, PHP, PKR, PLN, PYG, QAR, RON, RSD, RUB, RWF, SAR, SBD, SCR,
	SDG, SEK, SGD, SHP, SLL, SOS, SRD, STD, SVC, SYP, SZL, THB, TJS, TMT,
	TND, TOP, TRY, TTD, TWD, TZS, UAH, UGX, USD, UYU, UZS, VEF, VND, VUV,
	WST, XAF, XAG, XAU, XCD, XDR, XOF, XPF, YER, ZAR, ZMK, ZMW, ZWL,
}

// NewAPILayerProvider inits APILayer provider record
func NewAPILayerProvider(options *rates.Options) *APILayer {
	apiLayer := new(APILayer)
	// init all units
	if len(options.Currencies) == 0 {
		options.Currencies = append(options.Currencies, APILayerCurrencies...)
	}
	for _, unit := range options.Currencies {
		if c, err := currency.ParseISO(unit); err == nil {
			apiLayer.currencies = append(apiLayer.currencies, c)
		}
	}
	if options.Token != "" {
		apiLayer.lastURL = apiLayerLastURL + "?access_key=" + options.Token
		if v, ok := options.Settings["date"]; ok {
			if date, ok := v.(string); ok {
				apiLayer.historyURL = apiLayerHistoryURL + "?access_key=" + options.Token + "&date=" + date
				apiLayer.date = date
			}
		} else {
			apiLayer.historyURL = apiLayer.lastURL
		}
		if len(options.Currencies) < 10 {
			currencies := strings.Join(options.Currencies, ",")
			apiLayer.lastURL += "&currencies=" + currencies
			apiLayer.historyURL += "&currencies=" + currencies
		}
	}

	return apiLayer
}

type apiLayerEnvelope struct {
	Success    bool
	Historical bool
	Terms      string
	Privacy    string
	TimeStamp  int64
	Source     string
	Quotes     map[string]interface{}
	Error      struct {
		Code uint16
		Type string
		Info string
	}
}

// Name returns name of the provider
func (al *APILayer) Name() string {
	return apiLayerName
}

// FetchLast gets exchange rates for the last day
func (al *APILayer) FetchLast() ([]rates.Rate, []error) {
	return al.fetch(al.lastURL)
}

// FetchHistory gets exchange rates for all existing days
func (al *APILayer) FetchHistory() ([]rates.Rate, []error) {
	return al.fetch(al.historyURL)
}

// FetchLast gets exchange rates for the last day
func (al *APILayer) fetch(url string) (alRates []rates.Rate, alErrors []error) {
	response, err := http.Get(url)
	if err != nil {
		alErrors = append(alErrors, err)
		return
	}
	defer response.Body.Close()

	var raw apiLayerEnvelope

	if err := json.NewDecoder(response.Body).Decode(&raw); err != nil {
		alErrors = append(alErrors, err)
		return
	}

	if !raw.Success {
		alErrors = append(alErrors, errors.New(raw.Error.Info))
		return
	}

	if raw.TimeStamp == 0 {
		alErrors = append(alErrors, errors.New("Could not determine date"))
		return
	}
	timestamp := time.Unix(raw.TimeStamp, 0).UTC()
	date := timestamp.Format(stdDateTime)
	for _, unit := range al.currencies {
		for item, value := range raw.Quotes {
			if item == USD+unit.String() {
				alRates = append(alRates, rates.Rate{
					Time:     timestamp,
					Date:     date,
					Base:     currency.USD,
					Unit:     unit,
					Currency: currency.USD.String() + "/" + unit.String(),
					Value:    value,
				})
			}
		}
	}
	return
}
