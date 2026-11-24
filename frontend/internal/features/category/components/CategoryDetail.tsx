type CategoryDetailProps = {
  slug: string
}

export function CategoryDetail({ slug }: CategoryDetailProps) {
  return (
    <>
      <h1>Category</h1>
      <p>{slug}</p>
    </>
  )
}
