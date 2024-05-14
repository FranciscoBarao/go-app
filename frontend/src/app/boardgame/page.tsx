import { GetAll } from '@/app/api/catalog/route'


export default async function Boardgame() {
  const data = await GetAll()
  return data
}