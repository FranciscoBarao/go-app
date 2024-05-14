import { Boardgame } from '@/app/api/catalog/model'

const baseURL = process.env.CATALOG_HOST as string

export async function Post(bg: Boardgame) {
    const res = await fetch(baseURL, {
        method: "POST",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(bg),
        cache: 'no-store'
    })

    if(res.ok){
        return res.json()
      }
    
    return <p>error</p>
}

export async function Get() {
    const res = await fetch(baseURL + '/1', {
        method: "GET",
        headers: {
            'Content-Type': 'application/json',
        },
    })
    const data = await res.json()
    console.log(data)
    console.log(data.publisher)
   // const boardgame = JSON.parse(data) 
    return data
}

export async function GetAll() {
    const res = await fetch(baseURL, {
        method: "GET",
        headers: {
            'Content-Type': 'application/json',
        },
        cache: 'no-store'
    })

    
    const data = await res.json()
    console.log(data)

    return (
    <div> 
        {data.map((bg: Boardgame) => {
            return (
                <div key={bg.name}>
                    <p>{bg.name}</p>
                    <p>{bg.publisher}</p>
                    <p>{bg.playerNumber}</p>
                </div>
            )
        })}
    </div>
    )
}

export async function Delete() {
    const res = await fetch(baseURL + '/1', {
        method: "DELETE",
        headers: {
            'Content-Type': 'application/json',
        },
    })
    const data = await res.json()
    return data
}

