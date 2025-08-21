package notifications

type Mailer interface {
	Send(templateFile NotificationMessageTemplate, recipients []string, subject string, data any) (int, error)
}

type NotificationService struct {
	Mailer Mailer
}

func NewNotificationService(n *NotificationService) *NotificationService {
	return &NotificationService{
		Mailer: n.Mailer,
	}
}
