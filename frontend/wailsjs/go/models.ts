export namespace main {
	
	export class KnobUpdate {
	    knob: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new KnobUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.knob = source["knob"];
	        this.value = source["value"];
	    }
	}

}

