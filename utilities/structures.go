package utilities

// Patient represents the patient information
type Patient struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	MobileNumber  string `json:"mobile_number"`
	DateTime      string `json:"date_time"`
	Comments      string `json:"comments"`
	ServiceNeeded string `json:"service_needed"`
}
