import {Category, Invoice} from "./models/models";

export function decodeHTML(str: string) : string {
    return new DOMParser().parseFromString(str, 'text/html').documentElement.textContent || ""; 
}

export function getCategoryByName(categories: Category[], name: string) : Category {
    for (let i = 0; i < categories.length; i++) {
      const element = categories[i];
      if (element.name == name) {
        return element;
      }
    }

    return new Category();
}

export function initInvoices(input: Invoice[]) {
    let data = [];
    input.forEach((inv) => {
        data.push(new Invoice(inv))
    })

    return data
}