import { createBrowserRouter } from 'react-router'

import { BoardgameListPage } from '../pages/BoardgameListPage'
import { BoardgamePage } from '../pages/BoardgamePage'
import { CategoryListPage } from '../pages/CategoryListPage'
import { CategoryPage } from '../pages/CategoryPage'
import { ContributorListPage } from '../pages/ContributorListPage'
import { ContributorPage } from '../pages/ContributorPage'
import { HomePage } from '../pages/HomePage'
import { MechanismListPage } from '../pages/MechanismListPage'
import { MechanismPage } from '../pages/MechanismPage'
import { Layout } from './Layout'

// Path table only. Each entry maps a URL to a page component.
export const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      { path: '/', element: <HomePage /> },
      { path: '/boardgames', element: <BoardgameListPage /> },
      { path: '/boardgames/:slug', element: <BoardgamePage /> },
      { path: '/categories', element: <CategoryListPage /> },
      { path: '/categories/:slug', element: <CategoryPage /> },
      { path: '/mechanisms', element: <MechanismListPage /> },
      { path: '/mechanisms/:slug', element: <MechanismPage /> },
      { path: '/contributors', element: <ContributorListPage /> },
      { path: '/contributors/:slug', element: <ContributorPage /> },
    ],
  },
])
