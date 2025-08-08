package outbox

const (
	TableName      = "outbox"
	IdempotencyKey = "idempotency_key"
	Status         = "status"
	Kind           = "kind"
	Data           = "data"
	CreatedAt      = "created_at"
	UpdatedAt      = "updated_at"
)
