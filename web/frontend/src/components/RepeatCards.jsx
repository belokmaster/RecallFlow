import { useState, useEffect } from 'react'
import './RepeatCards.css'

function RepeatCards({ onBack }) {
  const [cards, setCards] = useState([])
  const [currentIndex, setCurrentIndex] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [feedback, setFeedback] = useState(null)

  useEffect(() => {
    loadCardsForReview()
  }, [])

  const loadCardsForReview = async () => {
    try {
      setLoading(true)
      setError(null)
      const res = await fetch('/repeat')
      const data = await res.json()
      if (res.ok) {
        setCards(data.cards || [])
        setCurrentIndex(0)
      } else {
        setError(data.error || 'Ошибка при загрузке карточек')
      }
    } catch (e) {
      setError('Ошибка сети')
      console.error(e)
    } finally {
      setLoading(false)
    }
  }

  const handleRepeat = async (success) => {
    if (currentIndex >= cards.length) return

    try {
      const cardId = cards[currentIndex].id
      const res = await fetch(`/repeat/${cardId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ success })
      })

      if (res.ok) {
        setFeedback(success ? 'Верно!' : 'Неправильно')
        setTimeout(() => {
          setFeedback(null)
          setCurrentIndex(currentIndex + 1)
        }, 800)
      } else {
        alert('Ошибка при сохранении результата')
      }
    } catch (e) {
      alert('Ошибка сети')
      console.error(e)
    }
  }

  if (loading) {
    return (
      <div className="repeat-container">
        <div className="loading">Загрузка карточек...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="repeat-container">
        <div className="error">{error}</div>
        <button className="btn-back" onClick={onBack}>Назад</button>
      </div>
    )
  }

  if (cards.length === 0) {
    return (
      <div className="repeat-container">
        <div className="empty-state">
          <div className="empty-icon">...</div>
          <h3>Все карточки повторены!</h3>
          <p>Отличная работа! Все доступные карточки повторены.</p>
          <button className="btn-action" onClick={onBack}>Вернуться к карточкам</button>
        </div>
      </div>
    )
  }

  if (currentIndex >= cards.length) {
    return (
      <div className="repeat-container">
        <div className="empty-state">
          <div className="empty-icon">...</div>
          <h3>Сессия завершена!</h3>
          <p>Отличная работа! Вы повторили {cards.length} карточек.</p>
          <button className="btn-action" onClick={onBack}>Вернуться к карточкам</button>
        </div>
      </div>
    )
  }

  const card = cards[currentIndex]
  const progress = `${currentIndex + 1} / ${cards.length}`

  return (
    <div className="repeat-container">
      <div className="repeat-header">
        <button className="btn-back" onClick={onBack}>Назад</button>
        <div className="progress-bar">
          <div className="progress-fill" style={{ width: `${((currentIndex + 1) / cards.length) * 100}%` }}></div>
        </div>
        <span className="progress-text">{progress}</span>
      </div>

      <div className="repeat-card">
        <div className="card-content">
          <div className="card-field">
            <label>Слово</label>
            <p className="card-word">{card.word}</p>
          </div>

          <div className="card-field">
            <label>Перевод</label>
            <p className="card-translation">{card.translation}</p>
          </div>

          {card.example && (
            <div className="card-field">
              <label>Пример</label>
              <p className="card-example">{card.example}</p>
            </div>
          )}

          <div className="card-stats">
            <span className="stat">Попыток: <strong>{card.attempts}</strong></span>
            <span className="stat">Верно: <strong>{card.successes}</strong></span>
          </div>
        </div>

        {feedback && (
          <div className={`feedback ${feedback.includes('✓') ? 'success' : 'error'}`}>
            {feedback}
          </div>
        )}
      </div>

      <div className="repeat-actions">
        <button 
          className="btn-success"
          onClick={() => handleRepeat(true)}
          disabled={feedback !== null}
        >
          ✓ Знаю
        </button>
        <button 
          className="btn-danger"
          onClick={() => handleRepeat(false)}
          disabled={feedback !== null}
        >
          ✗ Не знаю
        </button>
      </div>
    </div>
  )
}

export default RepeatCards
