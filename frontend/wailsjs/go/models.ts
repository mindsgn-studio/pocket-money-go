export namespace database {
	
	export class Wallet {
	    uuid: string;
	    name: string;
	    type: string;
	    address: string;
	
	    static createFrom(source: any = {}) {
	        return new Wallet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.address = source["address"];
	    }
	}

}

