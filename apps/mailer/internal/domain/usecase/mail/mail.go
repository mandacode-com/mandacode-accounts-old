package maildomain

type MailUsecase interface {
	// SendEmailVerificationMail sends an email verification mail to the user.
	//
	// Parameters:
	//   - email: The email address of the user to send the verification mail to.
	//   - verificationLink: The link to be included in the email for verification.
	SendEmailVerificationMail(email string, verificationLink string) error
}
