package tokens

import (
    "errors"
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

func (tm *TokenManager) CreateToken(userID, value int) (*models.StatusToken, error) {
    token := &models.StatusToken{
        UserID:    userID,
        Value:     value,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour), // Token expires after 24 hours
    }

    err := tm.DB.CreateToken(token)
    if err != nil {
        return nil, err
    }

    return token, nil
}

func (tm *TokenManager) UseToken(tokenID, userID int) error {
    token, err := tm.DB.GetTokenByID(tokenID)
    if err != nil {
        return err
    }

    if !token.IsValid() {
        return errors.New("token has expired")
    }

    if token.UserID != userID {
        return errors.New("token does not belong to this user")
    }

    err = tm.DB.MarkTokenAsUsed(tokenID)
    if err != nil {
        return err
    }

    return nil
}

func (tm *TokenManager) TransferToken(tokenID, currentOwnerID, newOwnerID int) error {
    token, err := tm.DB.GetTokenByID(tokenID)
    if err != nil {
        return err
    }

    if !token.IsValid() {
        return errors.New("token has expired")
    }

    if token.UserID != currentOwnerID {
        return errors.New("token does not belong to the current owner")
    }

    err = tm.DB.UpdateTokenOwner(tokenID, newOwnerID)
    if err != nil {
        return err
    }

    return nil
}
