package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog/log"

	"qrcodegen/internal/dto"
	"qrcodegen/internal/repository/postgres"
	sqldb "qrcodegen/sqlc/generated"

	"github.com/jackc/pgx/v5"
	"github.com/ua-parser/uap-go/uaparser"
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

type GeoResolver interface {
	Resolve(ip string) (string, string, bool)
}

type LinkUseCase struct {
	repo     postgres.Repository
	uaParser *uaparser.Parser
	geo      GeoResolver
}

func NewLinkUseCase(repo postgres.Repository, geo GeoResolver) *LinkUseCase {
	parser := uaparser.NewFromSaved()
	return &LinkUseCase{repo: repo, uaParser: parser, geo: geo}
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
		Name:        req.Name,
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
		Name:        linkData.Name,
		Color:       linkData.Color,
		Background:  linkData.Background,
		Smoothing:   linkData.Smoothing,
	}

	return response, nil
}

func (uc *LinkUseCase) GetAllLinks(ctx context.Context, userID int64) (*dto.GetAllLinksResponse, error) {
	links, err := uc.repo.GetLinksSummaryByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.GetAllLinksResponse{
				Links:   []dto.LinkInfo{},
				Message: "Success get all links by user",
			}, nil
		}
		return nil, fmt.Errorf("failed to get links by user id: %w", err)
	}

	linkInfos := make([]dto.LinkInfo, len(links))
	for i, link := range links {
		transitionsCount, _ := link.TransitionsCount.(int64)
		linkInfos[i] = dto.LinkInfo{ID: link.ID, OriginalURL: link.OriginalUrl, Name: link.Name, CreatedAt: link.CreatedAt, Transitions: transitionsCount}
	}

	return &dto.GetAllLinksResponse{
		Links:   linkInfos,
		Message: "Success get all links by user",
	}, nil
}

func (uc *LinkUseCase) SearchLinksByName(ctx context.Context, userID int64, search string) (*dto.GetAllLinksResponse, error) {
	params := sqldb.SearchLinksSummaryByNameParams{
		UserID:  userID,
		Column2: &search,
	}
	rows, err := uc.repo.SearchLinksSummaryByName(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &dto.GetAllLinksResponse{Links: []dto.LinkInfo{}, Message: "Success get all links by user"}, nil
		}
		return nil, fmt.Errorf("failed to search links by name: %w", err)
	}

	linkInfos := make([]dto.LinkInfo, len(rows))
	for i, link := range rows {
		transitionsCount, _ := link.TransitionsCount.(int64)
		linkInfos[i] = dto.LinkInfo{ID: link.ID, OriginalURL: link.OriginalUrl, Name: link.Name, CreatedAt: link.CreatedAt, Transitions: transitionsCount}
	}
	return &dto.GetAllLinksResponse{Links: linkInfos, Message: "Success get all links by user"}, nil
}

type LinkSortBy string

const (
	SortByCreatedAt   LinkSortBy = "created_at"
	SortByTransitions LinkSortBy = "transitions"
)

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

func (uc *LinkUseCase) SortLinks(items []dto.LinkInfo, by LinkSortBy, order SortOrder) []dto.LinkInfo {
	sorter := func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) }
	switch by {
	case SortByTransitions:
		sorter = func(i, j int) bool {
			if order == SortAsc {
				return items[i].Transitions < items[j].Transitions
			}
			return items[i].Transitions > items[j].Transitions
		}
	case SortByCreatedAt:
		fallthrough
	default:
		sorter = func(i, j int) bool {
			if order == SortAsc {
				return items[i].CreatedAt.Before(items[j].CreatedAt)
			}
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
	}
	sort.Slice(items, sorter)
	return items
}

func (uc *LinkUseCase) EditLink(ctx context.Context, linkID int64, userID int64, req dto.EditLinkRequest) (*dto.EditLinkResponse, error) {
	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	repoWithTx := uc.repo.WithTX(tx)

	updateLinkParams := sqldb.UpdateLinkURLParams{
		OriginalUrl: req.OriginalURL,
		ID:          linkID,
		UserID:      userID,
	}
	tag, err := repoWithTx.UpdateLinkURL(ctx, updateLinkParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update link: %w", err)
	}
	if tag == 0 {
		return nil, ErrLinkNotFound
	}

	updateQRParams := sqldb.UpdateQRCodeParamsParams{
		Color:      req.Color,
		Background: req.Background,
		Smoothing:  &req.Smoothing,
		LinkID:     linkID,
	}
	err = repoWithTx.UpdateQRCodeParams(ctx, updateQRParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update qr code params: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &dto.EditLinkResponse{
		Message: "Link updated successfully",
		ID:      linkID,
	}, nil
}

func (uc *LinkUseCase) Redirect(ctx context.Context, hash, referer, userAgent, ip string) (string, error) {
	link, err := uc.repo.GetLinkByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrLinkNotFound
		}
		return "", fmt.Errorf("failed to get link by hash: %w", err)
	}

	go func() {
		ctxBg, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		uc.createTransition(ctxBg, link.ID, referer, userAgent, ip)
	}()

	return link.OriginalUrl, nil
}

func (uc *LinkUseCase) createTransition(ctx context.Context, linkID int64, referer, userAgent, ip string) {
	var refPtr, uaPtr *string
	if referer != "" {
		refPtr = &referer
	}
	if userAgent != "" {
		uaPtr = &userAgent
	}

	client := uc.uaParser.Parse(userAgent)
	var brPtr, osPtr *string
	if client.UserAgent.Family != "" {
		brPtr = &client.UserAgent.Family
	}
	if client.Os.Family != "" {
		osPtr = &client.Os.Family
	}

	var countryPtr, cityPtr *string
	if uc.geo != nil && ip != "" {
		if country, city, ok := uc.geo.Resolve(ip); ok {
			if country != "" {
				countryPtr = &country
			}
			if city != "" {
				cityPtr = &city
			}
		}
	}

	params := sqldb.CreateTransitionParams{
		LinkID:    linkID,
		Country:   countryPtr,
		City:      cityPtr,
		Referer:   refPtr,
		UserAgent: uaPtr,
		Browser:   brPtr,
		Os:        osPtr,
	}

	err := uc.repo.CreateTransition(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create transition")
	}
}

func (uc *LinkUseCase) GetTransitions(ctx context.Context, linkID, userID int64) (*dto.GetTransitionsResponse, error) {
	type row = struct {
		ID        int64
		Country   *string
		City      *string
		Referer   *string
		UserAgent *string
		Browser   *string
		Os        *string
		CreatedAt time.Time
	}

	rows, err := uc.repo.GetTransitionsByLinkID(
		ctx,
		sqldb.GetTransitionsByLinkIDParams{
			LinkID: linkID,
			UserID: userID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transitions: %w", err)
	}

	items := make([]dto.TransitionItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.TransitionItem{
			ID:        r.ID,
			Country:   r.Country,
			City:      r.City,
			Referer:   r.Referer,
			UserAgent: r.UserAgent,
			Browser:   r.Browser,
			OS:        r.Os,
			CreatedAt: r.CreatedAt,
		})
	}

	return &dto.GetTransitionsResponse{Transitions: items}, nil
}

func (uc *LinkUseCase) DeleteLink(ctx context.Context, linkID int64, userID int64) error {
	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	repoWithTx := uc.repo.WithTX(tx)

	_, err = repoWithTx.GetLinkAndQRCodeByID(ctx, sqldb.GetLinkAndQRCodeByIDParams{ID: linkID, UserID: userID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrLinkNotFound
		}
		return fmt.Errorf("failed to verify link ownership: %w", err)
	}

	if err := repoWithTx.DeleteTransitionsByLinkID(ctx, linkID); err != nil {
		return fmt.Errorf("failed to delete transitions: %w", err)
	}

	if err := repoWithTx.DeleteQRCodeByLinkID(ctx, linkID); err != nil {
		return fmt.Errorf("failed to delete qr code: %w", err)
	}

	rowsAffected, err := repoWithTx.DeleteLink(ctx, sqldb.DeleteLinkParams{ID: linkID, UserID: userID})
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}
	if rowsAffected == 0 {
		return ErrLinkNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
