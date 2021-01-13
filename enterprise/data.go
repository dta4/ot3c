package enterprise

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

//ConsumptionAPIResponds represents the different Ressources that can be queried over the Enterprice API
type ConsumptionAPIResponds struct {
	Offset   int           `json:"offset"`
	Total    int           `json:"total"`
	Limit    int           `json:"limit"`
	Kind     string        `json:"kind"`
	Contents []Consumption `json:"contents"`
}

//Consumption is a single Ressource from the enterprice dashboard api
type Consumption struct {
	Kind               string  `json:"kind"`
	Contract           int     `json:"contract"`
	BusinessPartnerID  int     `json:"business_partner_id"`
	ResellerID         int     `json:"reseller_id"`
	QuantityUnit       string  `json:"quantity_unit"`
	ConsumptionDate    string  `json:"consumption_date"`
	UnitPrice          float64 `json:"unit_price"`
	ListpriceAmount    float64 `json:"listprice_amount"`
	Amount             float64 `json:"amount"`
	Product            string  `json:"product"`
	ProductDescription string  `json:"product_description"`
	ProjectID          string  `json:"project_id"`
	ResourceID         string  `json:"resource_id"`
	BillingQuantity    float64 `json:"billing_quantity"`
	ProjectName        string  `json:"project_name"`
}

//ConsumptionBetween returns a Consumption Array with Consumption in between t1 and t2
func ConsumptionBetween(c []Consumption, t1 time.Time, t2 time.Time) []Consumption {

	var result []Consumption
	for _, con := range c {
		t, err := time.Parse("2006-01-02", strings.Split(con.ConsumptionDate, " ")[0])
		if err != nil {
			logrus.WithError(err).Error("ay")
		}

		if t.After(t1) && t.Before(t2) {
			result = append(result, con)
		}
	}
	return result
}
