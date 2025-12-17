package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oriiyx/fritz/app/core/utils/env"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

type Manager struct {
	config  *env.Conf
	queries *db.Queries
	logger  *zerolog.Logger
}

func NewManager(config *env.Conf, queries *db.Queries, logger *zerolog.Logger) *Manager {
	l := logger.With().Str("service", "manager").Logger()

	return &Manager{
		config:  config,
		queries: queries,
		logger:  &l,
	}
}

func (m *Manager) CreateSession(ctx context.Context, userID pgtype.UUID, deviceID string) (*db.Session, error) {
	expiresAt := time.Now().Add(m.config.Session.SessionDuration)

	// Create session using queries
	session, err := m.queries.CreateSession(ctx, db.CreateSessionParams{
		UserIdentityID: userID,
		DeviceID:       deviceID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &session, nil
}

func (m *Manager) ValidateSession(ctx context.Context, sessionID string) (*db.Session, error) {
	var sessionUUID pgtype.UUID
	err := sessionUUID.Scan(sessionID)
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to scan UUID trying to get SessionID while validating session")
		return nil, err
	}

	session, err := m.queries.GetSessionDetailsByID(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	if !session.IsActive {
		return nil, errors.New("session is not active")
	}

	if time.Now().UTC().After(session.ExpiresAt.Time.UTC()) {
		// Deactivate expired session
		err = m.queries.DeactivateSession(ctx, session.ID)
		if err != nil {
			m.logger.Error().Err(err).Msg("failed to deactivate expired session")
		}
		return nil, errors.New("session has expired")
	}

	// Update last activity
	_, err = m.queries.UpdateSessionActivity(ctx, db.UpdateSessionActivityParams{
		ID:       sessionUUID,
		IsActive: true,
	})
	if err != nil {
		m.logger.Error().Err(err).Msg("failed to update session activity")
		// Don't return error here as the session is still valid
	}

	return &session, nil
}

func (m *Manager) InvalidateSession(ctx context.Context, sessionID string) error {
	var sessionUUID pgtype.UUID
	err := sessionUUID.Scan(sessionID)
	if err != nil {
		m.logger.Error().
			Err(err).
			Str("session_id", sessionID).
			Msg("Failed to scan UUID trying to get SessionID while validating session")
		return err
	}

	return m.queries.DeactivateSession(ctx, sessionUUID)
}

func (m *Manager) StartCleanup(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := m.queries.CleanupExpiredSessions(ctx); err != nil {
					m.logger.Error().Err(err).Msg("failed to cleanup expired sessions")
				}
			}
		}
	}()
}
