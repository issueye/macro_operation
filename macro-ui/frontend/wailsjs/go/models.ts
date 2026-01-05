export namespace main {
	
	export class MacroInfo {
	    Name: string;
	    CreatedAt: string;
	    UpdatedAt: string;
	    Script: string;
	
	    static createFrom(source: any = {}) {
	        return new MacroInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.CreatedAt = source["CreatedAt"];
	        this.UpdatedAt = source["UpdatedAt"];
	        this.Script = source["Script"];
	    }
	}

}

