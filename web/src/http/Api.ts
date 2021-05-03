const backEndHost: string =
    (process.env.NODE_ENV === 'production') ? 'https://api.bscfees.com' : 'http://localhost:3000'

export function api<T>(path: string): Promise<T> {
    return fetch(`${backEndHost}/${path}`, {mode: 'cors'})
        .then(
            response => {
                if (!response.ok) {
                    throw new Error(response.statusText)
                }
                return response.json() as Promise<T>
            })
        .then(data => data)
}