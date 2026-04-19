package infrastructure

import (
	"errors"
	"sync"

	"github.com/edmarqueslima/memorygame/internal/domain"
)

var ErrGameNotFound = errors.New("jogo não encontrado")

type MemoryRepo struct {
	mu    sync.RWMutex
	games map[string]*domain.Game
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		games: make(map[string]*domain.Game),
	}
}

func (r *MemoryRepo) Save(game *domain.Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.games[game.ID] = game
	return nil
}

func (r *MemoryRepo) Get(id string) (*domain.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	game, exists := r.games[id]
	if !exists {
		return nil, ErrGameNotFound
	}
	return game, nil
}
