import { Post } from '@/app/api/catalog/route'
import { Boardgame } from '@/app/api/catalog/model'


export default async function CreateBoardgame() {
  // TODO: Data from form as props
  let bg = new Boardgame("name","pubs", 3)
  
  // TODO: Error handling
  const data = await Post(bg)
  return (
    <div>
      <p>{data.name}</p>
      <p>{data.publisher}</p>
    </div>
  )
}