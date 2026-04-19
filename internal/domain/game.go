package domain

import (
	"errors"
	"math/rand"
	"time"
)

var (
	ErrGameOver     = errors.New("o jogo já acabou")
	ErrCardNotFound = errors.New("carta não encontrada")
	ErrCardFlipped  = errors.New("carta já está virada")
	ErrTwoCardsOpen = errors.New("duas cartas já estão viradas, verifique o par primeiro")
)

var defaultEmojis = []string{"🍎", "🍌", "🍒", "🍇", "🍉", "🍓", "🥑", "🍍"}

type Card struct {
	ID        int    `json:"id"`
	Emoji     string `json:"emoji,omitempty"` // Empty if face down
	IsFlipped bool   `json:"is_flipped"`
	IsMatched bool   `json:"is_matched"`
	realEmoji string // Hidden from JSON
}

type Game struct {
	ID           string  `json:"id"`
	Cards        []*Card `json:"cards"`
	Moves        int     `json:"moves"`
	Status       string  `json:"status"` // "playing", "won"
	FlippedIndex []int   `json:"-"`      // Track currently flipped cards
}

func NewGame(id string) *Game {
	cards := make([]*Card, 0, len(defaultEmojis)*2)
	idCounter := 0

	// Create pairs
	for _, emoji := range defaultEmojis {
		cards = append(cards, &Card{ID: idCounter, realEmoji: emoji})
		idCounter++
		cards = append(cards, &Card{ID: idCounter, realEmoji: emoji})
		idCounter++
	}

	// Shuffle
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
		// Reassign IDs after shuffle so they represent the board position
		cards[i].ID = i
		cards[j].ID = j
	})

	return &Game{
		ID:           id,
		Cards:        cards,
		Moves:        0,
		Status:       "playing",
		FlippedIndex: []int{},
	}
}

// FlipCard tries to flip a card. If 2 are already flipped, it expects CheckMatch to be called.
func (g *Game) FlipCard(index int) error {
	if g.Status != "playing" {
		return ErrGameOver
	}
	if index < 0 || index >= len(g.Cards) {
		return ErrCardNotFound
	}
	card := g.Cards[index]
	if card.IsFlipped || card.IsMatched {
		return ErrCardFlipped
	}

	// If 2 cards are already flipped, we must check them first
	// but the UI should have called CheckMatch. We'll do it automatically here if needed.
	if len(g.FlippedIndex) == 2 {
		g.CheckMatch()
	}

	card.IsFlipped = true
	card.Emoji = card.realEmoji // Reveal to UI
	g.FlippedIndex = append(g.FlippedIndex, index)

	return nil
}

// CheckMatch checks if currently flipped cards match.
func (g *Game) CheckMatch() {
	if len(g.FlippedIndex) != 2 {
		return
	}
	g.Moves++
	idx1 := g.FlippedIndex[0]
	idx2 := g.FlippedIndex[1]

	if g.Cards[idx1].realEmoji == g.Cards[idx2].realEmoji {
		g.Cards[idx1].IsMatched = true
		g.Cards[idx2].IsMatched = true
	} else {
		g.Cards[idx1].IsFlipped = false
		g.Cards[idx1].Emoji = ""
		g.Cards[idx2].IsFlipped = false
		g.Cards[idx2].Emoji = ""
	}
	g.FlippedIndex = []int{}
	g.checkWin()
}

func (g *Game) checkWin() {
	for _, c := range g.Cards {
		if !c.IsMatched {
			return
		}
	}
	g.Status = "won"
}

// HideEmojis hides emojis of unflipped cards before sending to frontend.
func (g *Game) SanitizeForClient() *Game {
	// We operate on a shallow copy or just modify the slice if it's for JSON serialization
	// Actually, the easiest is to just clear Emoji on unflipped cards.
	for _, c := range g.Cards {
		if !c.IsFlipped && !c.IsMatched {
			c.Emoji = ""
		} else {
			c.Emoji = c.realEmoji
		}
	}
	return g
}
