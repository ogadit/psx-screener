export namespace scraper {
	
	export class LineItem {
	    label: string;
	    isPercent: boolean;
	    values: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new LineItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.isPercent = source["isPercent"];
	        this.values = source["values"];
	    }
	}
	export class FinancialData {
	    companyName: string;
	    currentPrice: number;
	    incomeStatement: LineItem[];
	    balanceSheet: LineItem[];
	    ratios: LineItem[];
	
	    static createFrom(source: any = {}) {
	        return new FinancialData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.companyName = source["companyName"];
	        this.currentPrice = source["currentPrice"];
	        this.incomeStatement = this.convertValues(source["incomeStatement"], LineItem);
	        this.balanceSheet = this.convertValues(source["balanceSheet"], LineItem);
	        this.ratios = this.convertValues(source["ratios"], LineItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

