package services

import (
	"fmt"
	"sync"
)

type NotifierGroup struct {
	Notifiers []Notifier
}

func NewNotifierGroup(notifiers ...Notifier) *NotifierGroup {
	return &NotifierGroup{
		Notifiers: notifiers,
	}
}

func (s *NotifierGroup) NotifySubscribers(
	creatorID uint64,
	notification *Notification,
) error {

	// Create the wait group for all the notifiers
	var wg sync.WaitGroup
	wg.Add(len(s.Notifiers))

	// Loop through all the notifiers
	for index := range s.Notifiers {
		go func(i int) {

			// Defer the cleanup function
			defer wg.Done()

			// Send the notification on the notifier
			if err := s.Notifiers[i].NotifySubscribers(creatorID, notification); err != nil {
				fmt.Println("Error notifying subscribers: ", err.Error())
			}

		}(index)
	}

	// Wait for them all to complete
	wg.Wait()

	// Return without error
	return nil

}
