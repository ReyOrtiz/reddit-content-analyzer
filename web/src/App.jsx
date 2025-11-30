import React, { useState } from 'react'
import ReactMarkdown from 'react-markdown'
import './App.css'
import SubredditsList from './components/SubredditsList'
import { searchRedditPosts } from './services/api'

function App() {
  const [searchMethod, setSearchMethod] = useState('search')
  const [topic, setTopic] = useState('')
  const [subreddits, setSubreddits] = useState(['golang'])
  const [limit, setLimit] = useState(1)
  const [threshold, setThreshold] = useState(0.5)
  const [createdAfter, setCreatedAfter] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [results, setResults] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    setResults(null)

    try {
      const response = await searchRedditPosts({
        topic,
        subreddits,
        limit,
        relevance_threshold: threshold,
        created_after: createdAfter || null,
        search_method: searchMethod,
      })
      setResults(response)
    } catch (err) {
      setError(err.message || 'An error occurred while searching Reddit posts')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Reddit Content Analyzer</h1>
        <p>Search for relevant subreddit posts by topic</p>
      </header>

      <form className="search-form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Search Method *</label>
          <div className="radio-group">
            <label className="radio-option">
              <input
                type="radio"
                name="searchMethod"
                value="search"
                checked={searchMethod === 'search'}
                onChange={(e) => setSearchMethod(e.target.value)}
              />
              <span>Search by Topic</span>
            </label>
            <label className="radio-option">
              <input
                type="radio"
                name="searchMethod"
                value="latest"
                checked={searchMethod === 'latest'}
                onChange={(e) => setSearchMethod(e.target.value)}
              />
              <span>Latest Posts</span>
            </label>
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="topic">Topic *</label>
          <input
            type="text"
            id="topic"
            value={topic}
            onChange={(e) => setTopic(e.target.value)}
            placeholder="e.g., CLI application development"
            required
          />
          {searchMethod === 'latest' && (
            <small className="field-hint">
              Topic is used for relevance evaluation even when fetching latest posts
            </small>
          )}
        </div>

        <div className="form-group">
          <label>Subreddits *</label>
          <SubredditsList
            subreddits={subreddits}
            onChange={setSubreddits}
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label htmlFor="limit">Limit Posts per Subreddit</label>
            <input
              type="number"
              id="limit"
              value={limit}
              onChange={(e) => setLimit(parseInt(e.target.value) || 25)}
              min="1"
              max="100"
            />
          </div>

          <div className="form-group">
            <label htmlFor="threshold">Relevance Threshold</label>
            <input
              type="number"
              id="threshold"
              value={threshold}
              onChange={(e) => setThreshold(parseFloat(e.target.value) || 0.5)}
              min="0"
              max="1"
              step="0.1"
            />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="createdAfter">Created After (Optional)</label>
          <input
            type="datetime-local"
            id="createdAfter"
            value={createdAfter}
            onChange={(e) => setCreatedAfter(e.target.value)}
          />
        </div>

        <button
          type="submit"
          className="submit-button"
          disabled={loading || !topic || subreddits.length === 0}
        >
          {loading ? 'Searching...' : 'Search Posts'}
        </button>
      </form>

      {error && <div className="error">{error}</div>}

      {loading && <div className="loading">Searching for posts...</div>}

      {results && (
        <div className="results">
          <div className="results-header">
            <h2>Search Results</h2>
            <p>
              Found {results.posts?.length || 0} posts matching your criteria
            </p>
          </div>
          {results.posts && results.posts.length > 0 ? (
            <div className="posts-list">
              {results.posts.map((post, index) => (
                <div key={index} className="post-card">
                  <div className="post-header">
                    <h3 className="post-title">
                      <a
                        href={post.url}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {post.title}
                      </a>
                    </h3>
                    <div className="post-header-right">
                      {post.is_relevant !== undefined && (
                        <span className={`relevance-indicator ${post.is_relevant ? 'relevant' : 'not-relevant'}`}>
                          {post.is_relevant ? '✓ Relevant' : '✗ Not Relevant'}
                        </span>
                      )}
                      <span className="subreddit-name">r/{post.subreddit_name}</span>
                    </div>
                  </div>
                  {post.content && (
                    <div className="post-content">
                      <div className="post-content-preview">
                        <ReactMarkdown>{post.content}</ReactMarkdown>
                      </div>
                    </div>
                  )}
                  {post.relevance_summary && (
                    <div className="relevance-summary">
                      <div className="relevance-summary-header">
                        <svg className="sparks-logo" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                          <path d="M6 2L7.5 8.5L14 10L7.5 11.5L6 18L4.5 11.5L-2 10L4.5 8.5L6 2Z" fill="#ff4500"/>
                          <path d="M18 12L18.5 14.5L21 15L18.5 15.5L18 18L17.5 15.5L15 15L17.5 14.5L18 12Z" fill="#ff4500"/>
                        </svg>
                        <strong>Relevance</strong>
                      </div>
                      <p>{post.relevance_summary}</p>
                    </div>
                  )}
                  <div className="post-meta">
                    <span>Score: {post.score}</span>
                    <span>Comments: {post.num_comments}</span>
                    <span>
                      Created: {new Date(post.created_at).toLocaleString()}
                    </span>
                    {post.relevance_score !== undefined && (
                      <span className="relevance-score-badge">
                        Relevance Score: <strong>{(post.relevance_score * 100).toFixed(1)}%</strong>
                      </span>
                    )}
                    <a
                      href={post.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="reddit-link"
                    >
                      View on Reddit →
                    </a>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p>No posts found matching your criteria.</p>
          )}
        </div>
      )}
    </div>
  )
}

export default App

