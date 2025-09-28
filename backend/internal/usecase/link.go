package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"qrcodegen/internal/dto"
	"qrcodegen/internal/repository/postgres"
	sqldb "qrcodegen/sqlc/generated"

	"github.com/jackc/pgx/v5"
)

const (
	defaultQRColor      = "000000"
	defaultQRBackground = "FFFFFF"
	hashLength          = 7
)

var (
	defaultQRSmoothing = 0.0
	ErrLinkNotFound    = errors.New("link not found or access denied")
)

type LinkUseCase struct {
	repo postgres.Repository
}

func NewLinkUseCase(repo postgres.Repository) *LinkUseCase {
	return &LinkUseCase{repo: repo}
}

func generateHash(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func (uc *LinkUseCase) CreateLink(ctx context.Context, req dto.CreateLinkRequest, userID int64) (*dto.CreateLinkResponse, error) {
	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	repoWithTx := uc.repo.WithTX(tx)

	var linkHash string
	for i := 0; i < 5; i++ {
		hash, err := generateHash(hashLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate hash: %w", err)
		}
		_, err = repoWithTx.GetLinkByHash(ctx, hash)
		if errors.Is(err, pgx.ErrNoRows) {
			linkHash = hash
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to check hash uniqueness: %w", err)
		}
	}
	if linkHash == "" {
		return nil, errors.New("could not generate a unique hash")
	}

	linkParams := sqldb.CreateLinkParams{
		OriginalUrl: req.OriginalURL,
		Hash:        linkHash,
		UserID:      userID,
	}
	createdLink, err := repoWithTx.CreateLink(ctx, linkParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}

	qrParams := sqldb.CreateQRCodeParams{
		LinkID:     createdLink.ID,
		Color:      defaultQRColor,
		Background: defaultQRBackground,
		Smoothing:  &defaultQRSmoothing,
	}
	if _, err = repoWithTx.CreateQRCode(ctx, qrParams); err != nil {
		return nil, fmt.Errorf("failed to create qr code: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &dto.CreateLinkResponse{
		ID:      createdLink.ID,
		Message: "Link created successfully",
	}, nil
}

func (uc *LinkUseCase) GetLinkByID(ctx context.Context, linkID int64, userID int64) (*dto.GetLinkResponse, error) {
	params := sqldb.GetLinkAndQRCodeByIDParams{
		ID:     linkID,
		UserID: userID,
	}
	linkData, err := uc.repo.GetLinkAndQRCodeByID(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLinkNotFound
		}
		return nil, fmt.Errorf("failed to get link by id: %w", err)
	}

	response := &dto.GetLinkResponse{
		ID:          linkData.ID,
		OriginalURL: linkData.OriginalUrl,
		Hash:        linkData.Hash,
		CreatedAt:   linkData.CreatedAt,
		UpdatedAt:   linkData.UpdatedAt,
		Color:       linkData.Color,
		Background:  linkData.Background,
		Smoothing:   linkData.Smoothing,
	}

	return response, nil
}
