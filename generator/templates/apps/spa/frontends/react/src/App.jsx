import { useEffect, useState } from 'react'

export default function App() {
  const [info, setInfo] = useState(null)
  const [error, setError] = useState(null)

  useEffect(() => {
    fetch('/api/version')
      .then((res) => res.json())
      .then(setInfo)
      .catch((err) => setError(String(err)))
  }, [])

  return (
    <main className="app">
      <h1>Goose + React</h1>
      <p>
        This page is served by the Goose SPA platform. The backend answers
        JSON under <code>/api</code>.
      </p>
      {info && (
        <p className="status">
          Backend says: <strong>{info.name}</strong> v{info.version}
        </p>
      )}
      {error && <p className="error">Backend unreachable: {error}</p>}
    </main>
  )
}
