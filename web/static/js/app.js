const gameBoard = document.getElementById('gameBoard');
const movesCount = document.getElementById('movesCount');
const btnRestart = document.getElementById('btnRestart');
const winMessage = document.getElementById('winMessage');

let currentGameId = null;
let isAnimating = false; // Prevent clicks while checking matches

async function startNewGame() {
    try {
        const response = await fetch('/api/games', { method: 'POST' });
        const game = await response.json();
        currentGameId = game.id;
        renderBoard(game);
        updateStats(game.moves);
        winMessage.classList.add('hidden');
    } catch (error) {
        console.error("Error starting new game:", error);
    }
}

function renderBoard(game) {
    gameBoard.innerHTML = '';
    game.cards.forEach(card => {
        const cardEl = document.createElement('div');
        cardEl.className = `card ${card.is_flipped ? 'flipped' : ''} ${card.is_matched ? 'matched' : ''}`;
        cardEl.dataset.index = card.id;

        // Front (Logically the back of the card, what user sees initially)
        const frontEl = document.createElement('div');
        frontEl.className = 'card-front';
        // You could put a logo or pattern here
        
        // Back (The side with the emoji)
        const backEl = document.createElement('div');
        backEl.className = 'card-back';
        backEl.textContent = card.emoji || ''; // Show emoji if available

        cardEl.appendChild(frontEl);
        cardEl.appendChild(backEl);

        cardEl.addEventListener('click', () => handleCardClick(card.id));
        gameBoard.appendChild(cardEl);
    });
}

async function handleCardClick(index) {
    if (isAnimating || !currentGameId) return;

    const cardEl = gameBoard.children[index];
    if (cardEl.classList.contains('flipped') || cardEl.classList.contains('matched')) {
        return;
    }

    try {
        const response = await fetch(`/api/games/${currentGameId}/flip`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ card_index: index })
        });
        
        if (!response.ok) {
            console.error(await response.text());
            return;
        }

        const game = await response.json();
        
        // Update the specific card to show flip
        cardEl.classList.add('flipped');
        cardEl.querySelector('.card-back').textContent = game.cards[index].emoji;

        // If two cards are now flipped, we need to check match
        // Count flipped but not matched cards in UI state
        const flippedCards = Array.from(gameBoard.children).filter(c => c.classList.contains('flipped') && !c.classList.contains('matched'));
        
        if (flippedCards.length === 2) {
            isAnimating = true; // Lock UI
            
            // Wait a moment so user can see the second card
            setTimeout(async () => {
                const checkRes = await fetch(`/api/games/${currentGameId}/check`, { method: 'POST' });
                const updatedGame = await checkRes.json();
                
                renderBoard(updatedGame);
                updateStats(updatedGame.moves);
                
                if (updatedGame.status === 'won') {
                    winMessage.classList.remove('hidden');
                }
                
                isAnimating = false; // Unlock UI
            }, 1000); // 1 second delay
        }

    } catch (error) {
        console.error("Error flipping card:", error);
    }
}

function updateStats(moves) {
    movesCount.textContent = moves;
}

btnRestart.addEventListener('click', startNewGame);

// Start game on load
document.addEventListener('DOMContentLoaded', startNewGame);
