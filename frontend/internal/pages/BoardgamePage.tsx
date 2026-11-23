import { useParams } from 'react-router'

import { BoardgameDetail } from '../features/boardgame'

export function BoardgamePage() {
  const { slug = '' } = useParams()
  return <BoardgameDetail slug={slug} />
}
