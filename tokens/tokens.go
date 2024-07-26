package tokens

import (
	"errors"
	"math/rand"
	"time"

	"github.com/maxfinnsjo/subad/database"
	"github.com/maxfinnsjo/subad/models"
)

type TokenManager struct {
	DB *database.DB
}

func NewTokenManager(db *database.DB) *TokenManager {
	return &TokenManager{DB: db}
}

func (tm *TokenManager) GenerateToken(userID int) (*models.StatusToken, error) {
	token := &models.StatusToken{
		UserID:    userID,
		Value:     rand.Intn(100) + 1, // Generate a random value between 1 and 100
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Token expires after 24 hours
	}

	err := tm.DB.CreateStatusToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (tm *TokenManager) GetUserTokens(userID int) ([]models.StatusToken, error) {
	return tm.DB.GetStatusTokensByUserID(userID)
}

func (tm *TokenManager) UseToken(tokenID int) error {
	return tm.DB.MarkTokenAsUsed(tokenID)
}

func (tm *TokenManager) CalculateUserStatus(userID int) (int, error) {
	tokens, err := tm.GetUserTokens(userID)
	if err != nil {
		return 0, err
	}

	totalValue := 0
	for _, token := range tokens {
		if token.IsValid() {
			totalValue += token.Value
		}
	}

	return totalValue, nil
}

func (tm *TokenManager) TradeToken(senderID, recipientID, tokenID int) error {
	token, err := tm.DB.GetTokenByID(tokenID)
	if err != nil {
		return err
	}

	if token.UserID != senderID {
		return errors.New("token does not belong to sender")
	}

	return tm.DB.UpdateTokenOwner(tokenID, recipientID)
}
