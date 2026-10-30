package point

import "go-gateway/internal/datalink/schema"

func applyPointUpdate(point *schema.Point, req UpdatePointRequest) error {
	if req.Name != nil {
		point.Name = *req.Name
	}
	if req.Description != nil {
		point.Description = *req.Description
	}
	if req.Address != nil {
		address, err := validateAndNormalizeAddress(*req.Address)
		if err != nil {
			return err
		}
		point.Address = address
	}
	if req.Function != nil {
		function, err := validateAndNormalizeFunction(*req.Function)
		if err != nil {
			return err
		}
		point.Function = function
	}
	if req.DataType != nil {
		if err := validatePointDataType(*req.DataType); err != nil {
			return err
		}
		point.DataType = *req.DataType
	}
	if req.DataFormat != nil {
		if err := validatePointDataFormat(*req.DataFormat); err != nil {
			return err
		}
		point.DataFormat = normalizePointDataFormat(*req.DataFormat)
	}
	if req.Mode != nil {
		point.Mode = *req.Mode
	}
	if req.PollingGroupID != nil || req.ReplacePollingGroup {
		point.PollingGroupID = req.PollingGroupID
	}
	if req.Enabled != nil {
		point.Enabled = *req.Enabled
	}

	return nil
}
