package application

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/edmarqueslima/memorygame/internal/domain"
)

type GameRepository interface {
	Save(game *domain.Game) error
	Get(id string) (*domain.Game, error)
}

type GameService struct {
	repo GameRepository
}

func NewGameService(repo GameRepository) *GameService {
	return &GameService{repo: repo}
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *GameService) StartNewGame() (*domain.Game, error) {
	id := generateID()
	game := domain.NewGame(id)
	err := s.repo.Save(game)
	if err != nil {
		return nil, err
	}
	return game.SanitizeForClient(), nil
}

func (s *GameService) GetGame(id string) (*domain.Game, error) {
	game, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return game.SanitizeForClient(), nil
}

func (s *GameService) FlipCard(gameID string, cardIndex int) (*domain.Game, error) {
	game, err := s.repo.Get(gameID)
	if err != nil {
		return nil, err
	}

	err = game.FlipCard(cardIndex)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(game)
	if err != nil {
		return nil, err
	}

	return game.SanitizeForClient(), nil
}

func (s *GameService) CheckMatch(gameID string) (*domain.Game, error) {
	game, err := s.repo.Get(gameID)
	if err != nil {
		return nil, err
	}

	game.CheckMatch()

	err = s.repo.Save(game)
	if err != nil {
		return nil, err
	}

	return game.SanitizeForClient(), nil
}
