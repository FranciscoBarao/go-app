type Boardgame = {
    name:         string    
	publisher:    string     
	playerNumber: number        
	tags:         Tag[]       
	categories:   Category[]
	mechanisms:   Mechanism[]
	expansions:   Boardgame[]
	boardgameID:  number       
}

type Tag = {
    name: string
}

type Mechanism = {
    name: string
}

type Category = {
    name: string
}