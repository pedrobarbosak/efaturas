<script lang="ts">
    import { onMount } from "svelte";
    import { translate } from "@/translations";
    import { logoutAndClear, user } from "@/store";

    import { Alert, Button, DarkMode, Spinner } from 'flowbite-svelte';
    import { BanOutline, InfoCircleSolid } from "flowbite-svelte-icons";

    import { Invoice } from "@/models/models";
    import Table from "@/components/pages/Table.svelte";
    import CategoryIcon from "@/components/pages/CategoryIcon.svelte";
    import {fetchNewInvoices, processInvoices, updateInvoices} from "@/api";

    export let data;
    let { invoices } = data;
    const { categories } = data;

    let lastDate = "---"
    let nToProcess = countToProcess();
    let nSelectedCategories = countSelectedCategories();

    const actions = { isLoading: false, text: "", error: "" };

    onMount(async () => {
        updateCategoriesTotal()

        if (invoices?.length === 0 ) {
            await actionFetch();
        }
    })

    $: if (invoices) {
        updateCategoriesTotal()
        nToProcess = countToProcess();
        nSelectedCategories = countSelectedCategories();

        if (invoices.length !== 0) {
            lastDate = invoices.reduce((p, c) => {return p?.document?.date > c?.document?.date ? p : c;})?.document.date || "---";
        }
    }

    function updateCategoriesTotal() {
        categories.forEach((cat) => { cat.total = 0; })

        invoices.forEach((inv) => {
            categories.forEach((cat) => {
                if (inv.activity.category == cat.name) {
                    cat.total += inv.total.benefit;
                }
            })
        })
    }


    function logout() {
        logoutAndClear()
        location.href = "/"
    }

    function countSelectedCategories(): number {
        let found = 0;

        for (let i=0; i<invoices.length; i++) {
            if (invoices[i].selected != "" && invoices[i].selected != invoices[i].activity.category) {
                found += 1;
            }
        }

        return found
    }

    function countToProcess(): number {
        let found = 0;

        for (let i=0; i<invoices.length; i++) {
            if (invoices[i].categories.size !== categories.length) {
                found += 1;
            }
        }

        return found
    }

    async function actionFetch() {
        actions.isLoading = true
        actions.text = translate("label_action_fetch_description")

        const resp = await fetchNewInvoices($user.token)
        if (resp.hasError) {
            actions.text = "";
            actions.isLoading = false;
            actions.error = resp.error || "Something went wrong"
            return
        }

        setInvoices(resp.data || [], false)
        clearActions()
    }


    async function actionProcess() {
        actions.isLoading = true
        actions.text = translate("label_action_process_description")

        const resp = await processInvoices($user.token)
        if (resp.hasError) {
            actions.text = "";
            actions.isLoading = false;
            actions.error = resp.error || "Something went wrong"
            return
        }

        setInvoices(resp.data || [], false)
        clearActions()
    }

    async function actionUpdateSelected() {
        actions.isLoading = true
        actions.text = translate("label_action_update_selected_description")

        const resp = await updateInvoices($user.token, invoices)
        if (resp.hasError) {
            actions.text = "";
            actions.isLoading = false;
            actions.error = resp.error || "Something went wrong"
            return
        }

        setInvoices(resp.data || [], false)
        clearActions()
    }

    function clearActions() {
        actions.text = ""
        actions.isLoading = false
    }

    function setInvoices(input: Invoice[], append: boolean = false) {
        if (!append) {
            invoices = [];
        }

        const invs: Invoice[] = [];
        input.forEach((inv) => {
            invs.push(new Invoice(inv))
        })

        if (append) {
            invoices.push(...invs)
            return
        }

        invoices = invs
    }
</script>


<!-- Form-->
<div class="flex items-center justify-between p-2">
    <div class="flex items-center w-1/4">
        <img src="/favicon.png" alt="logo" class="w-24"/>
    </div>

    <div class="flex items-center justify-center gap-2 w-2/4">
        <Button outline pill disabled={actions.isLoading} on:click={actionFetch}>
            { translate("label_action_fetch") }
        </Button>

        <Button outline pill disabled={actions.isLoading || nToProcess === 0} on:click={actionProcess}>
            { translate("label_action_process") } ({ nToProcess })
        </Button>

        <Button outline pill disabled={actions.isLoading || nSelectedCategories === 0} on:click={actionUpdateSelected}>
            { translate("label_action_update_selected") } ({ nSelectedCategories })
        </Button>
    </div>

    <div class="flex items-center justify-end gap-2 w-1/4">
        <Button color="alternative" class="border-gray-600" outline pill on:click={()=>{ logout() }}>
            <BanOutline class="w-5 h-5 me-2" />{translate("label_logout")}
        </Button>
        <DarkMode class="border border-gray-600 rounded-3xl text-black" />
    </div>
</div>


<!-- Status -->
<div class="flex flex-col items-center justify-center p-4 border-t border-b rounded">
    {#if actions.text}
        <div class="dark:text-white">
            <Spinner /> { actions.text }
        </div>
    {:else if actions.error}
        <Alert color="red" dismissable>
            <InfoCircleSolid slot="icon" class="w-5 h-5" />
            <span class="font-medium">Error!</span>
            { actions.error }
        </Alert>
    {:else}
        <div class="dark:text-white">
            <strong>{ translate("label_last_invoice_at") }:</strong>
            { lastDate }
        </div>
    {/if}

    <span class="text-sm text-gray-400">
        Invoices: { invoices.length }
    </span>
</div>


<!-- Categories -->
<div class="flex items-center justify-center  rounded">
    {#each categories as category}
        <CategoryIcon showTotal={true} {category}/>
    {/each}
</div>


<!-- Table -->
<Table {categories} bind:invoices={invoices} />
