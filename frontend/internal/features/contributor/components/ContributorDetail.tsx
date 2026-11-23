type ContributorDetailProps = {
  slug: string
}

export function ContributorDetail({ slug }: ContributorDetailProps) {
  return (
    <>
      <h1>Contributor</h1>
      <p>{slug}</p>
    </>
  )
}
