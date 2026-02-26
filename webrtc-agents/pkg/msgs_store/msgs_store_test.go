package msgs_store

import (
	"sync"
	"testing"
)

// MockMessage is a mock message for testing
type MockMessage struct {
	Sender string
	Seq    int
}

// GetSenderId implements IdentifiableMessage interface
func (m *MockMessage) GetSenderId() string {
	return m.Sender
}

// Sender wraps SyncMsgsStore and sends messages with a specific sender name
type Sender struct {
	store  *SyncMsgsStore
	sender string
}

// NewSender creates a new Sender
func NewSender(store *SyncMsgsStore, sender string) *Sender {
	return &Sender{
		store:  store,
		sender: sender,
	}
}

// SendMessages sends N fake messages to the store
func (s *Sender) SendMessages(n int) error {
	for i := 0; i < n; i++ {
		msg := &MockMessage{
			Sender: s.sender,
			Seq:    i,
		}
		if err := s.store.Append(msg); err != nil {
			return err
		}
	}
	return nil
}

func TestConcurrentSendMessages(t *testing.T) {
	store := NewSyncMsgsStore()

	// Define test cases with different senders and message counts
	testCases := []struct {
		sender string
		count  int
	}{
		{"sender1", 100},
		{"sender2", 150},
		{"sender3", 200},
		{"sender4", 50},
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(testCases))

	for _, tc := range testCases {
		wg.Add(1)
		go func(sender string, count int) {
			defer wg.Done()
			s := NewSender(store, sender)
			if err := s.SendMessages(count); err != nil {
				errChan <- err
			}
		}(tc.sender, tc.count)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			t.Errorf("Error sending messages: %v", err)
		}
	}

	// Verify the store contains all messages
	result := store.Load()
	if result == nil {
		t.Fatal("Store is nil after sending messages")
	}

	// Count total messages expected
	totalExpected := 0
	for _, tc := range testCases {
		totalExpected += tc.count
	}

	// Verify we have the correct number of messages
	coll := result.Load().(*IndexedMsgsCollection)
	totalActual := 0
	for _, messages := range coll.store {
		totalActual += len(messages)
	}

	if totalActual != totalExpected {
		t.Errorf("Expected %d messages, got %d", totalExpected, totalActual)
	}

	// Verify each sender has the correct number of messages
	for _, tc := range testCases {
		messages := coll.store[tc.sender]
		if len(messages) != tc.count {
			t.Errorf("Sender %s: expected %d messages, got %d", tc.sender, tc.count, len(messages))
		}
	}
}

func TestSenderSendMessages(t *testing.T) {
	store := NewSyncMsgsStore()
	sender := NewSender(store, "test-sender")

	err := sender.SendMessages(10)
	if err != nil {
		t.Fatalf("SendMessages returned error: %v", err)
	}

	result := store.Load()
	if result == nil {
		t.Fatal("Store is nil after sending messages")
	}

	coll := result.Load().(*IndexedMsgsCollection)
	messages := coll.store["test-sender"]

	if len(messages) != 10 {
		t.Errorf("Expected 10 messages, got %d", len(messages))
	}

	// Verify sequence numbers
	for i, msg := range messages {
		mockMsg := msg.(*MockMessage)
		if mockMsg.Seq != i {
			t.Errorf("Message %d: expected seq %d, got %d", i, i, mockMsg.Seq)
		}
		if mockMsg.Sender != "test-sender" {
			t.Errorf("Message %d: expected sender 'test-sender', got '%s'", i, mockMsg.Sender)
		}
	}
}
