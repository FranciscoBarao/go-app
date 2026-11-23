import { useParams } from 'react-router'

import { ContributorDetail } from '../features/contributor'

export function ContributorPage() {
  const { slug = '' } = useParams()
  return <ContributorDetail slug={slug} />
}
