package services

import "sync"

type LiveChatMessageBuffer struct {
	MaxLength int
	items     []*ChatMsg
	mut       sync.RWMutex
}

func (buf *LiveChatMessageBuffer) Push(msg *ChatMsg) {

	// Lock on the mutex with write access
	buf.mut.Lock()
	defer buf.mut.Unlock()

	// Add the message to the buffer
	buf.items = append(buf.items, msg)

	// Count how many items we have in excess
	excess := len(buf.items) - buf.MaxLength

	// If there is an excess of items
	if excess > 0 {

		// Create a new buffer for the items
		newItems := make([]*ChatMsg, buf.MaxLength)

		// Loop through them and copy over
		for i := excess; i < len(buf.items); i++ {
			newItems[i-excess] = buf.items[i]
		}

		// Assign the new slice
		buf.items = newItems

	}

}

func (buf *LiveChatMessageBuffer) GetCopy() []*ChatMsg {

	// Lock on the mutex with readonly access
	buf.mut.RLock()
	defer buf.mut.RUnlock()

	// Create the new slice for elements
	items := make([]*ChatMsg, len(buf.items))

	// Copy all the elements
	for i := range buf.items {
		items[i] = buf.items[i]
	}

	// Return the new slice
	return items

}
