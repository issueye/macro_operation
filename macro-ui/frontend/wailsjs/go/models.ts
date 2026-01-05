export namespace main {
	
	export class EventInfo {
	    index: number;
	    type: string;
	    keyCode: number;
	    chars: string;
	    x: number;
	    y: number;
	    button: number;
	    delta: number;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new EventInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.type = source["type"];
	        this.keyCode = source["keyCode"];
	        this.chars = source["chars"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.button = source["button"];
	        this.delta = source["delta"];
	        this.timestamp = source["timestamp"];
	    }
	}
	export class MacroInfo {
	    name: string;
	    createdAt: string;
	    updatedAt: string;
	    script: string;
	    event_count: number;
	
	    static createFrom(source: any = {}) {
	        return new MacroInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.script = source["script"];
	        this.event_count = source["event_count"];
	    }
	}

}

