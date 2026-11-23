// Application settings. Vite inlines VITE_* values at build time, so a change needs a new build.
const DEFAULT_CATALOG_BASE_URL = 'http://localhost:8081'

export const config = {
  mode: import.meta.env.MODE,
  baseUrl: import.meta.env.BASE_URL,
  catalog: {
    baseUrl: import.meta.env.VITE_CATALOG_BASE_URL || DEFAULT_CATALOG_BASE_URL,
    // TODO: request timeout once the catalog client sends requests.
  },
  pagination: {
    defaultPage: 1,
    defaultPageSize: 10,
  },
}
