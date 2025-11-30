import React, { useState } from 'react'
import './SubredditsList.css'

function SubredditsList({ subreddits, onChange }) {
  const [newSubreddit, setNewSubreddit] = useState('')
  const [editingIndex, setEditingIndex] = useState(null)
  const [editValue, setEditValue] = useState('')

  const handleAdd = () => {
    const trimmed = newSubreddit.trim().toLowerCase()
    if (trimmed && !subreddits.includes(trimmed)) {
      onChange([...subreddits, trimmed])
      setNewSubreddit('')
    }
  }

  const handleRemove = (index) => {
    onChange(subreddits.filter((_, i) => i !== index))
  }

  const handleEditStart = (index) => {
    setEditingIndex(index)
    setEditValue(subreddits[index])
  }

  const handleEditSave = (index) => {
    const trimmed = editValue.trim().toLowerCase()
    if (trimmed) {
      const updated = [...subreddits]
      updated[index] = trimmed
      onChange(updated)
    }
    setEditingIndex(null)
    setEditValue('')
  }

  const handleEditCancel = () => {
    setEditingIndex(null)
    setEditValue('')
  }

  const handleKeyPress = (e, action) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      action()
    }
  }

  return (
    <div className="subreddits-list">
      <div className="subreddits-input-group">
        <input
          type="text"
          value={newSubreddit}
          onChange={(e) => setNewSubreddit(e.target.value)}
          onKeyPress={(e) => handleKeyPress(e, handleAdd)}
          placeholder="Enter subreddit name (e.g., golang)"
          className="subreddit-input"
        />
        <button
          type="button"
          onClick={handleAdd}
          className="add-button"
          disabled={!newSubreddit.trim()}
        >
          Add
        </button>
      </div>

      <div className="subreddits-tags">
        {subreddits.map((subreddit, index) => (
          <div key={index} className="subreddit-tag">
            {editingIndex === index ? (
              <div className="subreddit-edit">
                <input
                  type="text"
                  value={editValue}
                  onChange={(e) => setEditValue(e.target.value)}
                  onKeyPress={(e) => handleKeyPress(e, () => handleEditSave(index))}
                  onBlur={() => handleEditSave(index)}
                  className="subreddit-edit-input"
                  autoFocus
                />
                <button
                  type="button"
                  onClick={() => handleEditSave(index)}
                  className="save-button"
                >
                  ✓
                </button>
                <button
                  type="button"
                  onClick={handleEditCancel}
                  className="cancel-button"
                >
                  ✕
                </button>
              </div>
            ) : (
              <>
                <span className="subreddit-name">r/{subreddit}</span>
                <button
                  type="button"
                  onClick={() => handleEditStart(index)}
                  className="edit-button"
                  title="Edit"
                >
                  ✎
                </button>
                <button
                  type="button"
                  onClick={() => handleRemove(index)}
                  className="remove-button"
                  title="Remove"
                >
                  ×
                </button>
              </>
            )}
          </div>
        ))}
      </div>

      {subreddits.length === 0 && (
        <p className="subreddits-empty">No subreddits added yet. Add at least one to search.</p>
      )}
    </div>
  )
}

export default SubredditsList

