// What the backend answers (backend/plugin.go).
export type Entry = { name: string; path: string; dir: boolean; size: number }
export type Summary = { folders: number; files: number; bytes: number }
