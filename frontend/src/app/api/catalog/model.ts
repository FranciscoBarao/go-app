export class Boardgame {
    name:         string    
	publisher:    string     
	playerNumber: number        
	tags:         Tag[]       
	categories:   string[]
	mechanisms:   string[]
	expansions:   Boardgame[]
	boardgameID?:  number       

	constructor(name: string, publisher: string, pnumber: number){
		this.name = name;
		this.publisher = publisher;
		this.playerNumber = pnumber;
		this.tags = [new Tag("B")]
		this.categories = []
		this.mechanisms = []
		this.expansions = []
	}
}

class Tag {
	name: string
	constructor(name: string){
		this.name = name
	}
}