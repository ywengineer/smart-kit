package mailx

type Mailer interface {
	// Send sends a mail to the specified recipient.
	Send(from, to, cc, bcc, subject, bodyType, bodyString string) error
}
