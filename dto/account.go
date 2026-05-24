package dto

// region Verify Code
type VerifyCodeResult struct {
	Token        string
	UserNotFound bool
	InvalidCode  bool
}
