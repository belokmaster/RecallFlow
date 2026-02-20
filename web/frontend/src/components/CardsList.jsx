import './CardsList.css'

export default function CardsList({ cards, onCardClick }) {
  const calculatePercent = (card) => {
    const attempts = card.attempts || 0
    if (attempts === 0) return 0
    return Math.round(((card.successes || 0) / attempts) * 100)
  }

  return (
    <div className="cards-grid">
      {cards.map(card => (
        <div 
          key={card.id} 
          className="card-item"
          onClick={() => onCardClick(card)}
        >
          <div className="card-word">{card.word}</div>
          <div className="card-translation">{card.translation}</div>
          <div className="card-stats-mini">
            <div className="stat-mini">
              <span>Попыток:</span>
              <span className="stat-value">{card.attempts || 0}</span>
            </div>
            <div className="stat-mini">
              <span>Успехов:</span>
              <span className="stat-value">{card.successes || 0}</span>
            </div>
            <div className="stat-mini">
              <span>%:</span>
              <span className="stat-value">{calculatePercent(card)}%</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
