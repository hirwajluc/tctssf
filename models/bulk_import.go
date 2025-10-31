package models

// CreateMemberRequest for admin to create members
type CreateMemberRequest struct {
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	MonthlyCommitment float64 `json:"monthly_commitment"`
	JoinedDate        string  `json:"joined_date,omitempty"`
}

// BulkMemberData represents a single member in bulk import
type BulkMemberData struct {
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	MonthlyCommitment float64 `json:"monthly_commitment"`
	Row               int     `json:"row"` // Track original row number for error reporting
}

// BulkImportResult represents the result of bulk import operation
type BulkImportResult struct {
	TotalProcessed int                    `json:"total_processed"`
	SuccessCount   int                    `json:"success_count"`
	ErrorCount     int                    `json:"error_count"`
	Errors         []BulkImportError      `json:"errors"`
	SuccessMembers []BulkImportSuccess    `json:"success_members"`
}

// BulkImportError represents an error during bulk import
type BulkImportError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// BulkImportSuccess represents a successful member creation
type BulkImportSuccess struct {
	Row           int    `json:"row"`
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
}