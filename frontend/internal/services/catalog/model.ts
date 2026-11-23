// Read shapes returned by the catalog API. Field names match the catalog JSON.

export type Page<T> = {
  data: T[]
  page: number
  pageSize: number
  totalItems: number
  totalPages: number
}

type Audited = {
  created_at: string
  updated_at: string
  deleted_at?: string
}

export type Category = Audited & {
  id: number
  slug: string
  name: string
  bgg_id?: number
}

export type Mechanism = Audited & {
  id: number
  slug: string
  name: string
  bgg_id?: number
}

export type Contributor = Audited & {
  id: number
  slug: string
  name: string
  bio?: string
  bgg_id?: number
}

export type ContributorRole =
  | 'designer'
  | 'artist'
  | 'publisher'
  | 'developer'
  | 'graphic_designer'

export type Contribution = {
  contributor: Contributor
  role: ContributorRole
  credit_order?: number
}

export type Boardgame = Audited & {
  id: number
  slug: string
  name: string
  description: string
  year_published: number
  min_players: number
  max_players: number
  min_play_time: number
  max_play_time: number
  min_age?: number
  bgg_id?: number
  boardgame_id?: number
  categories?: Category[]
  mechanisms?: Mechanism[]
  contributions?: Contribution[]
  expansions?: Boardgame[]
}
