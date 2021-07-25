package services

type Notification struct {
	Title string
	Body  string
	Link  *string
	Image *string
}

type Notifier interface {
	NotifySubscribers(creatorID uint64, notification *Notification) error
}
