package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/Timothylock/go-signin-with-apple/apple"
)

type AuthenticatedAppleUser struct {
	AppleUserId string
	Email       string
}

func ValidateAuthorizationToken(
	token string,
	privateKey string,
	clientID string,
	teamID string,
	keyID string,
) (*AuthenticatedAppleUser, error) {
	secret, err := apple.GenerateClientSecret(strings.ReplaceAll(privateKey, "@#", "\n"), teamID, clientID, keyID)
	fmt.Println("secret", secret, err)

	if err != nil {
		return nil, err
	}

	client := apple.New()

	req := apple.AppValidationTokenRequest{
		ClientID:     clientID,
		ClientSecret: secret,
		Code:         token,
	}

	var resp apple.ValidationResponse

	// Do the verification
	err = client.VerifyAppToken(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		fmt.Printf("apple returned an error: %s - %s\n", resp.Error, resp.ErrorDescription)
		if err != nil {
			return nil, err
		}
	}

	// Get the unique user ID
	userId, err := apple.GetUniqueID(resp.IDToken)
	if err != nil {
		return nil, err
	}

	// Get the email
	claim, err := apple.GetClaims(resp.IDToken)
	if err != nil {
		return nil, err
	}

	email := (*claim)["email"].(string)

	return &AuthenticatedAppleUser{
		AppleUserId: userId,
		Email:       strings.TrimSpace(strings.ToLower(email)),
	}, nil
}
