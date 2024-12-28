package feed

import "proj2/lock"

// Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	Add(body string, timestamp float64)
	Remove(timestamp float64) bool
	Contains(timestamp float64) bool
	GetAllPosts() []*post // Added: Method to get all posts
}

// feed is the internal representation of a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation. You can assume the feed will not have duplicate posts
type feed struct {
	start *post        // a pointer to the beginning post
	rw    *lock.RWLock // read/write lock for thread safety
}

// post is the internal representation of a post on a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation.
type post struct {
	body      string  // the text of the post
	timestamp float64 // Unix timestamp of the post
	next      *post   // the next post in the feed
}

// NewPost creates and returns a new post value given its body and timestamp
func newPost(body string, timestamp float64, next *post) *post {
	return &post{body, timestamp, next}
}

// NewFeed creates a empy user feed
func NewFeed() Feed {
	return &feed{start: nil, rw: lock.NewRWLock()}
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *feed) Add(body string, timestamp float64) {
	f.rw.Lock() // Acquire a write lock for thread-safe operations
	defer f.rw.Unlock() // Ensure the lock is released

	// Handle case when the feed is empty or the new post is the most recent
	if f.start == nil || f.start.timestamp < timestamp {
		// Insert at the beginning
		f.start = newPost(body, timestamp, f.start) // Create a new post as the start of the feed
		return
	}

	// Traverse the feed to find the appropriate position for the new post
	current := f.start
	for current.next != nil && current.next.timestamp >= timestamp {
		current = current.next
	}

	// Insert the new post
	current.next = newPost(body, timestamp, current.next)
}

// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *feed) Remove(timestamp float64) bool {
	f.rw.Lock() // Acquire a write lock for thread-safe operations
	defer f.rw.Unlock() // Ensure the lock is released at the end of the function

	// Handle case when the feed is empty
	if f.start == nil {
		return false // No posts to remove
	}

	// Check if the post to remove is at the start
	if f.start.timestamp == timestamp {
		f.start = f.start.next // Remove the start node
		return true
	}

	// Traverse the feed to find the post with the specified timestamp
	current := f.start
	for current.next != nil && current.next.timestamp != timestamp {
		current = current.next // Move to the next post
	}

	// If the post (timestamp) was not found, return false
	if current.next == nil {
		return false
	}

	// Remove the post
	current.next = current.next.next
	return true
}

// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *feed) Contains(timestamp float64) bool {
	f.rw.RLock()
	defer f.rw.RUnlock()

	// Traverse the feed to check for the post with the given timestamp
	current := f.start
	for current != nil {
		if current.timestamp == timestamp {
			return true // Post found
		}
		current = current.next
	}
	// Post not found
	return false
}

// Body returns the body of the post
func (p *post) Body() string {
	return p.body
}

// Timestamp returns the timestamp of the post
func (p *post) Timestamp() float64 {
	return p.timestamp
}


// Implement GetAllPosts for feed
func (f *feed) GetAllPosts() []*post {
	f.rw.RLock()
	defer f.rw.RUnlock()

	var posts []*post
	current := f.start
	for current != nil {
		posts = append(posts, current)
		current = current.next
	}
	return posts
}

