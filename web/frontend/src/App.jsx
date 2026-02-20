import { useState, useEffect } from 'react'
import './App.css'
import CardsList from './components/CardsList'
import CreateCardModal from './components/CreateCardModal'
import ViewCardModal from './components/ViewCardModal'
import RepeatCards from './components/RepeatCards'

function App() {
  const [cards, setCards] = useState([])
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showViewModal, setShowViewModal] = useState(false)
  const [selectedCard, setSelectedCard] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [currentPage, setCurrentPage] = useState('cards')

  useEffect(() => {
    loadCards()
  }, [])

  const loadCards = async () => {
    try {
      setLoading(true)
      setError(null)
      const res = await fetch('/cards')
      const data = await res.json()
      if (res.ok) {
        setCards(data.cards || [])
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

  const handleCreateCard = async (word, translation, example) => {
    try {
      const res = await fetch('/cards', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          word: word.trim(),
          translation: translation.trim(),
          example: example.trim() || null
        })
      })
      
      if (res.ok) {
        setShowCreateModal(false)
        loadCards()
      } else {
        const err = await res.json()
        alert(err.error || 'Ошибка при создании карточки')
      }
    } catch (e) {
      alert('Ошибка сети')
      console.error(e)
    }
  }

  const handleDeleteCard = async (id) => {
    if (!confirm('Удалить эту карточку?')) return
    
    try {
      const res = await fetch(`/cards/${id}`, { method: 'DELETE' })
      if (res.ok) {
        setShowViewModal(false)
        loadCards()
      } else {
        alert('Ошибка при удалении карточки')
      }
    } catch (e) {
      alert('Ошибка сети')
      console.error(e)
    }
  }

  const handleViewCard = (card) => {
    setSelectedCard(card)
    setShowViewModal(true)
  }

  return (
    <div className="app">
      {currentPage === 'repeat' ? (
        <RepeatCards onBack={() => setCurrentPage('cards')} />
      ) : (
        <>
      <header className="app-header">
        <div className="header-content">
          <h1 className="brand">Recall Flow</h1>
          <div className="header-buttons">
            <button 
              className="btn-repeat" 
              onClick={() => setCurrentPage('repeat')}
            >
              Повторение
            </button>
            <button 
              className="btn-create" 
              onClick={() => setShowCreateModal(true)}
            >
              <span>+</span> Новая карточка
            </button>
          </div>
        </div>
      </header>

      <main className="container">
        <div className="section-header">
          <h2 className="section-title">Мои карточки</h2>
          <span className="section-badge">{cards.length}</span>
        </div>

        {loading ? (
          <div className="loading">Загрузка...</div>
        ) : error ? (
          <div className="error">{error}</div>
        ) : cards.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">...</div>
            <h3>Нет карточек</h3>
            <p>Начните с создания первой карточки для изучения новых слов</p>
            <button className="btn-action" onClick={() => setShowCreateModal(true)}>
              Создать карточку
            </button>
          </div>
        ) : (
          <CardsList cards={cards} onCardClick={handleViewCard} />
        )}
      </main>

      {showCreateModal && (
        <CreateCardModal 
          onClose={() => setShowCreateModal(false)}
          onCreate={handleCreateCard}
        />
      )}

      {showViewModal && selectedCard && (
        <ViewCardModal 
          card={selectedCard}
          onClose={() => setShowViewModal(false)}
          onDelete={handleDeleteCard}
        />
      )}
        </>
      )}
    </div>
  )
}

export default App
