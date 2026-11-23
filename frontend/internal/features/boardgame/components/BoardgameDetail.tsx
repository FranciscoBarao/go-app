type BoardgameDetailProps = {
  slug: string
}

export function BoardgameDetail({ slug }: BoardgameDetailProps) {
  return (
    <>
      <h1>Boardgame</h1>
      <p>{slug}</p>
    </>
  )
}
