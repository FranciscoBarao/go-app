import { useParams } from 'react-router'

import { MechanismDetail } from '../features/mechanism'

export function MechanismPage() {
  const { slug = '' } = useParams()
  return <MechanismDetail slug={slug} />
}
