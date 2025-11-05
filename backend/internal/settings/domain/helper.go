package domain

// StringPtr returns a pointer to the string value
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// StringValue returns the string value or empty string if nil
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// EmiratePtr returns a pointer to the emirate value
func EmiratePtr(e string) *UAEmirate {
	if e == "" {
		return nil
	}
	emirate := UAEmirate(e)
	return &emirate
}

// EmirateValue returns the emirate value or empty string if nil
func EmirateValue(e *UAEmirate) string {
	if e == nil {
		return ""
	}
	return string(*e)
}
