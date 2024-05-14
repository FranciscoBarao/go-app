import { Get } from '@/app/api/catalog/route'


export default async function Boardgame() {
  const data = await Get()
  return data
}