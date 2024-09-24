package utilities

// Patient represents the patient information
type Patient struct {
	PatientName   string `json:"name"`
	EmailAddress  string `json:"email_address"`
	PhoneNumber   string `json:"phone_number"`
	City          string `json:"city"`
	DateTime      string `json:"date_time"`
	ServiceNeeded string `json:"service_needed"`
	Comments      string `json:"comments"`
}
