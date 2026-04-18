package infrastructure

import (
	"context"
	"errors"
	"os"

	"google.golang.org/api/idtoken"
)

// VerifyGoogleIDToken verifies the id_token from Google and returns the payload.
func VerifyGoogleIDToken(idToken string) (*idtoken.Payload, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID not set")
	}

	payload, err := idtoken.Validate(context.Background(), idToken, clientID)
	if err != nil {
		return nil, err
	}

	// Google guarantees email_verified=true for the accounts we care about
	if emailVerified, ok := payload.Claims["email_verified"].(bool); !ok || !emailVerified {
		return nil, errors.New("email not verified by Google")
	}

	return payload, nil
}
