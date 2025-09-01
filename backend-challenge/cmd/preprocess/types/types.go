package types

// Rec represents a hashed coupon record.
type Rec struct {
	H    uint64 // Hash of the coupon code
	Code string // Coupon code
}
