import { config } from '../../config/config'
import type { Boardgame, Category, Contributor, Mechanism, Page } from './model'

// Catalog client stubs. Every function returns an empty result until the HTTP calls are added.
//
// TODO: send requests to config.catalog.baseUrl. Paths follow catalog/internal/route:
//   boardgame    GET /api/boardgame, GET /api/boardgame/{slug}
//   category     GET /api/category, GET /api/category/{slug}
//   mechanism    GET /api/mechanism, GET /api/mechanism/{slug}
//   contributor  GET /api/contributors, GET /api/contributors/{slug}
// TODO: send page and pageSize as query params (?page=&pageSize=) on list calls.

function emptyPage<T>(page: number, pageSize: number): Page<T> {
  return { data: [], page, pageSize, totalItems: 0, totalPages: 0 }
}

export async function listBoardgames(
  page = config.pagination.defaultPage,
  pageSize = config.pagination.defaultPageSize,
): Promise<Page<Boardgame>> {
  return emptyPage(page, pageSize)
}

export async function getBoardgame(_slug: string): Promise<Boardgame | null> {
  return null
}

export async function listCategories(
  page = config.pagination.defaultPage,
  pageSize = config.pagination.defaultPageSize,
): Promise<Page<Category>> {
  return emptyPage(page, pageSize)
}

export async function getCategory(_slug: string): Promise<Category | null> {
  return null
}

export async function listMechanisms(
  page = config.pagination.defaultPage,
  pageSize = config.pagination.defaultPageSize,
): Promise<Page<Mechanism>> {
  return emptyPage(page, pageSize)
}

export async function getMechanism(_slug: string): Promise<Mechanism | null> {
  return null
}

export async function listContributors(
  page = config.pagination.defaultPage,
  pageSize = config.pagination.defaultPageSize,
): Promise<Page<Contributor>> {
  return emptyPage(page, pageSize)
}

export async function getContributor(
  _slug: string,
): Promise<Contributor | null> {
  return null
}
