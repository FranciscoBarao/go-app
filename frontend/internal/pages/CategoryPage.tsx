import { useParams } from 'react-router'

import { CategoryDetail } from '../features/category'

export function CategoryPage() {
  const { slug = '' } = useParams()
  return <CategoryDetail slug={slug} />
}
