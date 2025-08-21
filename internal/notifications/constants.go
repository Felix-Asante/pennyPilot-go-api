package notifications

type NotificationMessageTemplate string

const (
	ForgotPasswordMessageTemplate NotificationMessageTemplate = "forgot-password.html"
)

func (n NotificationMessageTemplate) String() string {
	return string(n)
}
