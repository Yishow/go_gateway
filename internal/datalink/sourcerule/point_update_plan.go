package sourcerule

import "go-gateway/internal/datalink/schema"

type pointUpdatePlan struct {
	appliedRevision string
	name            string
	address         string
	point           *schema.Point
	batch           bool
}
