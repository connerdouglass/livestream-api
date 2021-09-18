package services

// RtmpAuthService manages authentication for the RTMP server when it calls into this
// API. By design, the RTMP server (github.com/connerdouglass/livestream-rtmp) has a
// hard-coded passcode.
//
// For more security later, we can validate the caller IP address, passcode, and more.
// Remember that no matter what, the stream key must also be valid for the creator streaming.
type RtmpAuthService struct {
	RtmpServerPasscode string
}

// CheckPasscode checks if the provided passcode is valid
func (s *RtmpAuthService) CheckPasscode(passcode string) bool {
	return passcode == s.RtmpServerPasscode
}
