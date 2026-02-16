package health

const (
	StatusOK = "ok"
)

type Status struct {
	Status string `json:"status"`
	Uptime int    `json:"uptime"`
}

func NewStatus(uptime int) *Status {
	result := &Status{
		Status: StatusOK,
		Uptime: uptime,
	}

	return result
}
