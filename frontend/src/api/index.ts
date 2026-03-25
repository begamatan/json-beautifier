const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'

export interface ApiError {
  code: string
  message: string
  details?: string
}

export interface OutputResponse {
  result: string
}

export interface ValidateResponse {
  valid: boolean
  message: string
}

async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  const data = await res.json()

  if (!res.ok) {
    throw data as ApiError
  }

  return data as T
}

export function beautify(json: string, indent: 2 | 4): Promise<OutputResponse> {
  return post<OutputResponse>('/beautify', { json, indent })
}

export function minify(json: string): Promise<OutputResponse> {
  return post<OutputResponse>('/minify', { json })
}

export function validate(json: string): Promise<ValidateResponse> {
  return post<ValidateResponse>('/validate', { json })
}
