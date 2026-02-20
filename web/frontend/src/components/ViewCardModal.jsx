import './ViewCardModal.css'

export default function ViewCardModal({ card, onClose, onDelete }) {
  const calculatePercent = () => {
    const attempts = card.attempts || 0
    if (attempts === 0) return 0
    return Math.round(((card.successes || 0) / attempts) * 100)
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content card-view-modal" onClick={e => e.stopPropagation()}>
        <button className="close-btn" onClick={onClose}>×</button>
        
        <div className="card-view-header">
          <div className="word-display">{card.word}</div>
          <div className="translation-display">{card.translation}</div>
        </div>

        {card.example && (
          <div className="example-section">
            <label>Пример:</label>
            <p>{card.example}</p>
          </div>
        )}

        <div className="stats-section">
          <div className="stat-item">
            <span className="stat-label">Попыток</span>
            <span className="stat-value">{card.attempts || 0}</span>
          </div>
          <div className="stat-item">
            <span className="stat-label">Успехов</span>
            <span className="stat-value">{card.successes || 0}</span>
          </div>
          <div className="stat-item">
            <span className="stat-label">Процент</span>
            <span className="stat-value">{calculatePercent()}%</span>
          </div>
        </div>

        <div className="modal-actions">
          <button 
            className="btn-delete"
            onClick={() => onDelete(card.id)}
          >
            Удалить
          </button>
          <button 
            className="btn-cancel"
            onClick={onClose}
          >
            Закрыть
          </button>
        </div>
      </div>
    </div>
  )
}
