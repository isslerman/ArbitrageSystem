package ntfy

import (
	"log/slog"
	"net/http"
	"strings"
)

type Ntfy struct{}

// NewNtfy creates an instance of the NTFY Service
// More info at: https://docs.ntfy.sh
func NewNtfy() *Ntfy {
	return &Ntfy{}
}

func (n *Ntfy) SendMsg(t, msg string, urgent bool) bool {
	req, _ := http.NewRequest("POST", "https://ntfy.sh/loop_magneto",
		strings.NewReader(msg))
	req.Header.Set("Title", t)
	if urgent {
		req.Header.Set("Priority", "urgent")
	}
	req.Header.Set("Tags", "warning,skull")

	// others headers
	// 	strings.NewReader(`There's someone at the door. 🐶

	//	Please check if it's a good boy or a hooman.
	//
	// Doggies have been known to ring the doorbell.`))

	// link to go when clock the msg
	// req.Header.Set("Click", "https://home.nest.com/")

	// attach
	// req.Header.Set("Attach", "https://nest.com/view/yAxkasd.jpg")

	// creates a button with a click action
	// req.Header.Set("Actions", "http, Open door, https://api.nest.com/open/yAxkasd, clear=true")

	// not sure?
	//req.Header.Set("Email", "phil@example.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(
			"error sending ntfy msg",
			slog.String("status", err.Error()),
		)
		return false
	}

	if res.StatusCode != 200 {
		slog.Error(
			"error sending ntfy msg",
			slog.Int("status", res.StatusCode),
		)
		return false
	}
	return true
}
