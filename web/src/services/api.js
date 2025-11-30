import axios from 'axios'

const API_BASE_URL = '/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const searchRedditPosts = async (params) => {
  try {
    // Format the request according to the backend DTO
    const requestData = {
      topic: params.topic || '',
      subreddits: params.subreddits,
      relevance_threshold: params.relevance_threshold || 0.5,
      limit: params.limit || 25,
      created_after: params.created_after
        ? new Date(params.created_after).toISOString()
        : null,
      min_num_comments: 0, // Default value
      search_method: params.search_method || 'search',
    }

    const response = await api.post('/reddit/relevance/search', requestData)
    return response.data
  } catch (error) {
    if (error.response) {
      throw new Error(
        error.response.data?.error || 'Failed to search Reddit posts'
      )
    } else if (error.request) {
      throw new Error('No response from server. Is the backend running?')
    } else {
      throw new Error(error.message || 'An error occurred')
    }
  }
}

