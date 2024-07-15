package cron

import (
	"github.com/liteseed/edge/internal/database/schema"
)

func (crn *Cron) JobReleaseReward() {
	o, err := crn.database.GetOrders(&schema.Order{Status: schema.Confirmed})
	if err != nil {
		crn.logger.Error("failed to fetch reward orders", err)
		return
	}

	for _, order := range *o {
		err = crn.contract.Release(order.ID)
		if err != nil {
			crn.logger.Error("failed to release reward", err)
			continue
		}
		err = crn.database.UpdateOrder(order.ID, &schema.Order{Status: schema.Permanent})
		if err != nil {
			crn.logger.Error("failed to update database", err)
			continue
		}
	}

}
