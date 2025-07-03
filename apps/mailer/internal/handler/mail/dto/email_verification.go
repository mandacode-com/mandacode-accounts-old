package mailhandlerdto

type MailVerificationRequest struct {
	// Email is the email address of the user to send the verification mail to.
	Email string `json:"email"`
	// VerificationLink is the link to be included in the email for verification.
	VerificationLink string `json:"verification_link"`
}
