import {decodeHTML} from "@/utils";

export class Invoice {
    id: number;
    origin: Origin;
    issuer: Issuer;
    buyer: Buyer;
    document: Document;
    total: Total;
    activity: Activity;
    atCud: string;
    categories: Map<string, Values> = new Map<string, Values>;

    selected: string;

    constructor(input?: Invoice) {
        this.id = input?.id || 0;
        this.origin = input?.origin || new Origin();
        this.issuer = input?.issuer || new Issuer();
        this.buyer = input?.buyer || new Buyer();
        this.document = input?.document || new Document();
        this.total = input?.total || new Total();
        this.activity = input?.activity || new Activity();
        this.atCud = input?.atCud || "";

        this.categories = new Map()
        if (input?.categories) {
            this.categories = new Map(Object.entries(input?.categories));
        }

        this.selected = this.activity.category;
    }

    static search(invoice: Invoice, search: string): boolean {
        const s = search.toLowerCase().trim()
        if (s == "") {
            return true
        }

        return (
            invoice.document.date.toLowerCase().includes(s) ||
            decodeHTML(invoice.document.number).toLowerCase().includes(s) ||
            decodeHTML(invoice.document.description).toLowerCase().includes(s) ||
            decodeHTML(invoice.issuer.nif).toLowerCase().includes(s) ||
            decodeHTML(invoice.issuer.name).toLowerCase().includes(s)
        )
    }
}

export class Values {
    success: boolean = false;
    benefit: number = 0;
    others: number = 0;
}

export class Origin {
    value: string = "";
    description: string = "";
}

export class Issuer {
    nif: string = "";
    name: string = "";
}

export class Buyer {
    nif: string = "";
    name: string = "";
    country: string = "";
    nifInternational: string = "";
}

export class Document {
    type: string = "";
    description: string = "";
    number: string = "";
    hash: string = "";
    date: string = "";
}

export class Total {
    value: number = 0;
    taxable: number = 0;
    vat: number = 0;
    benefit: number = 0;
}

export class Activity {
    category: string = "";
    description: string = "";
}

export class Category {
    name: string = "";
    color: string = "";
    unicode: string = "";
    total: number = 0;
}

export class Money {
    static format(n: number) : string {
        return (n / 100).toFixed(2) + " €"
    }
}