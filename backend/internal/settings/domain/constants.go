package domain

// backend/settings/internal/domain/shafafiya.go

// Add these constants at the end of your shafafiya.go file:

const (
	EnvironmentProduction = "production"
	EnvironmentStaging    = "staging"
	EnvironmentTest       = "test"
)

// ---- Shafafiya Web Service Operations ----

// Transaction operations
const (
	OperationUploadTransaction        = "UploadTransaction"
	OperationGetNewTransactions       = "GetNewTransactions"
	OperationDownloadTransactionFile  = "DownloadTransactionFile"
	OperationSetTransactionDownloaded = "SetTransactionDownloaded"
	OperationSearchTransactions       = "SearchTransactions"
)

// DRG (Diagnosis Related Group) operations
const (
	OperationAddDRGToEClaim = "AddDRGToEClaim"
	OperationGetDRGDetails  = "GetDRGDetails"
)

// Prior Authorization operations
const (
	OperationCheckForNewPriorAuthorizationTransactions = "CheckForNewPriorAuthorizationTransactions"
	OperationGetNewPriorAuthorizationTransactions      = "GetNewPriorAuthorizationTransactions"
)

// Prescription operations
const (
	OperationGetPrescriptions = "GetPrescriptions"
)

// V3 Additional operations (not available in V2)
const (
	OperationGetPersonInsuranceHistory            = "GetPersonInsuranceHistory"
	OperationGetClaimCountReconciliation          = "GetClaimCountReconciliation"
	OperationGetInsuranceContinuityCertificate    = "GetInsuranceContinuityCertificate"
	OperationCancelInsuranceContinuityCertificate = "CancelInsuranceContinuityCertificate"
)

// ---- Submission Constants ----

const (
	SubmissionStatusSuccess = "SUCCESS"
	SubmissionStatusFailure = "FAILURE"
	SubmissionStatusPending = "PENDING"
)

// ---- Costing Method Constants ----

const (
	CostingMethodDepartmental  = "DEPARTMENTAL"
	CostingMethodActivityBased = "ACTIVITY_BASED"
	CostingMethodServiceBased  = "SERVICE_BASED"
)

// ---- Allocation Method Constants ----

const (
	AllocationMethodWeighted   = "WEIGHTED"
	AllocationMethodPercentage = "PERCENTAGE"
	AllocationMethodFixed      = "FIXED"
)

// ---- Language Constants ----

const (
	LanguageEnglish = "en"
	LanguageArabic  = "ar"
)

// ---- Currency Constants ----

const (
	CurrencyAED = "AED" // UAE Dirham
	CurrencyUSD = "USD" // US Dollar
	CurrencyEUR = "EUR" // Euro
	CurrencyGBP = "GBP" // British Pound
	CurrencyINR = "INR" // Indian Rupee
	CurrencySAR = "SAR" // Saudi Riyal
	CurrencyKWD = "KWD" // Kuwaiti Dinar
	CurrencyQAR = "QAR" // Qatari Riyal
)
