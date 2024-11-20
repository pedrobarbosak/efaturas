import { login } from "@/api";
import { Invoice } from "./models";
import {plainToClass} from "class-transformer";

export class UILoading<T> {
	value: T;
	isLoading: boolean;
	isDisabled: boolean;
	verified?: boolean;

	constructor(value: T, isDisabled?: boolean, isLoading?: boolean) {
		this.value = value;
		this.isDisabled = isDisabled ?? false;
		this.isLoading = isLoading ?? false;
	}
}

export class Form {
    username: UILoading<string> = new UILoading("");
    password: UILoading<string> = new UILoading("");
	remember: boolean = false
	isLoading: boolean = false
	error: string | null = null

    async onSubmit(): Promise<string> {
		this.isLoading = true;

		const resp = await login(this.username.value, this.password.value)
		if (resp.hasError) {
			this.error = resp.error
			this.isLoading = false
			return ""
		}

		return resp.data?.token || ""
    }

	async fetchStuff() {
		const resp = await login(this.username.value, this.password.value)
		if (resp.hasError) {
			this.error = resp.error
			this.isLoading = false
			return []
		}

		let invoices: Invoice[] = [];
		const data = resp.data as Invoice[];
		data.forEach((inv) => {
			invoices.push(new Invoice(inv))
		})

		// this.myMap = new Map(Object.entries(mapObj));
		// let invoices: Invoice[] = [];
		//
		// const inputs = resp.data as Invoice[]
		// inputs.forEach((inv) => {
		// 	invoices.push(new Invoice(inv))
		// })
		//
		//return resp.data as Invoice[]

		return invoices;
	}
}
