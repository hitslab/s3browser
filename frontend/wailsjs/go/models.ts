export namespace config {
	
	export class S3Settings {
	    base_url: string;
	    region: string;
	    bucket: string;
	    access_key: string;
	    secret_key: string;
	    path_style: boolean;
	    disable_ssl: boolean;
	
	    static createFrom(source: any = {}) {
	        return new S3Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.base_url = source["base_url"];
	        this.region = source["region"];
	        this.bucket = source["bucket"];
	        this.access_key = source["access_key"];
	        this.secret_key = source["secret_key"];
	        this.path_style = source["path_style"];
	        this.disable_ssl = source["disable_ssl"];
	    }
	}
	export class Config {
	    s3: S3Settings;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.s3 = this.convertValues(source["s3"], S3Settings);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace s3 {
	
	export class List {
	    folders: string[];
	    files: string[];
	
	    static createFrom(source: any = {}) {
	        return new List(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.folders = source["folders"];
	        this.files = source["files"];
	    }
	}

}

