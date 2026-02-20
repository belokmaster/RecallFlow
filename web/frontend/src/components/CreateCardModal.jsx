import { useState } from 'react'
import './CreateCardModal.css'

export default function CreateCardModal({ onClose, onCreate }) {
  const [word, setWord] = useState('')
  const [translation, setTranslation] = useState('')
  const [example, setExample] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    if (!word.trim() || !translation.trim()) {
      alert('Заполните слово и перевод')
      return
    }

    setLoading(true)
    try {
      await onCreate(word, translation, example)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <button className="close-btn" onClick={onClose}>×</button>
        
        <h3 className="modal-title">Новая карточка</h3>
        
        <form onSubmit={handleSubmit} className="create-form">
          <div className="form-group">
            <label htmlFor="word">Слово <span className="required">*</span></label>
            <input
              id="word"
              type="text"
              value={word}
              onChange={(e) => setWord(e.target.value)}
              placeholder="Например: Hello"
              maxLength="250"
              disabled={loading}
              autoFocus
            />
          </div>

          <div className="form-group">
            <label htmlFor="translation">Перевод <span className="required">*</span></label>
            <input
              id="translation"
              type="text"
              value={translation}
              onChange={(e) => setTranslation(e.target.value)}
              placeholder="Например: Привет"
              maxLength="250"
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="example">Пример (опционально)</label>
            <textarea
              id="example"
              value={example}
              onChange={(e) => setExample(e.target.value)}
              placeholder="Пример использования слова..."
              disabled={loading}
              rows="3"
            />
          </div>

          <button 
            type="submit" 
            className="btn-submit"
            disabled={loading}
          >
            {loading ? 'Создание...' : 'Создать'}
          </button>
        </form>
      </div>
    </div>
  )
}
