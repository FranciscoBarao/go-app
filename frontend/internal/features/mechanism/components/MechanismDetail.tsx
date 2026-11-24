type MechanismDetailProps = {
  slug: string
}

export function MechanismDetail({ slug }: MechanismDetailProps) {
  return (
    <>
      <h1>Mechanism</h1>
      <p>{slug}</p>
    </>
  )
}
