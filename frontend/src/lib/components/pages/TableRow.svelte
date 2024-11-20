<script lang="ts">
    import { Category, Invoice, Money } from '@/models/models';
    import { TableBodyCell, TableBodyRow } from 'flowbite-svelte';
    import CategoryIcon from '@/components/pages/CategoryIcon.svelte'
    import { decodeHTML, getCategoryByName } from '@/utils';
    import { onDestroy, onMount } from "svelte";
    import { getInvoiceSSE } from '@/api';
    import { fly, fade } from 'svelte/transition';
    import { Indicator } from 'flowbite-svelte';
    import {user} from "@/store";

    export let invoice: Invoice = new Invoice(); 
    export let categories: Category[] = [];
    export let isOpen: boolean = false;
    export let hide: boolean = false;

    let eventSource : EventSource;

    $: cats = Array.from(invoice.categories?.entries())

    $: category = getCategoryByName(categories, invoice.activity.category)

    onMount(async ()=> {
        if (invoice.id == 0 || invoice.id == undefined || invoice.categories.size == categories.length) {
            return
        }

        const onEvent = (event) => {
            const data: { done: boolean, invoice: Invoice } = JSON.parse(event.data)
            invoice = new Invoice(data.invoice);
        }

        await getInvoiceSSE($user.token, invoice.id, onEvent)
    })

    onDestroy(() => {
      if(eventSource) {
        eventSource.close()
      }
    })
</script>

{#if !hide}
    <TableBodyRow class="hover:border-1 hover:border-[{category.color}]" on:click={()=> {isOpen = !isOpen}}>
        <TableBodyCell class="truncate max-w-xs">
            <div class="flex flex-row items-center p-0">
                <Indicator color="yellow" size="xs" class="{invoice.selected === invoice.activity.category ? 'invisible opacity-0' : ''}"/>
                <CategoryIcon size={24} category={category} />
            </div>
        </TableBodyCell>

        <TableBodyCell class="truncate max-w-xs">
            {invoice.document.date}
        </TableBodyCell>

        <TableBodyCell class="truncate max-w-xs">
            {decodeHTML(invoice.document.number)} - {decodeHTML(invoice.document.description)}
        </TableBodyCell>

        <TableBodyCell class="truncate max-w-xs">
            <div title="{invoice.issuer.nif} - {decodeHTML(invoice.issuer.name)}">
                {invoice.issuer.nif} - {decodeHTML(invoice.issuer.name)}
            </div>
        </TableBodyCell>

        <TableBodyCell class="flex">
            {#each Array.from(invoice.categories?.entries()) as [category, values]}
                <div in:fly out:fade>
                    <CategoryIcon size={16} category={getCategoryByName(categories, category)} disabled={!values.success} />
                </div>
            {/each}
        </TableBodyCell>

        <TableBodyCell>
            {Money.format(invoice.total.value)}
        </TableBodyCell>

<!--        <TableBodyCell>-->
<!--            {Money.format(invoice.total.benefit)}-->
<!--        </TableBodyCell>-->

    </TableBodyRow>

    {#if isOpen}
        <TableBodyRow>
            <TableBodyCell colspan={6} class=""> <!-- 7 -->
                <div class="flex justify-center items-center pb-10 mr-2 border-b border-l border-r" in:fly out:fade style="z-index: 1000">
                    <div class="flex flex-row p-1">
                        {#if cats.length !== 0}
                            {#each cats as [category, values]}
                                {#if values.success}
                                    <div class="flex flex-row p-2 hover:scale-[1.15] rounded-3xl cursor-pointer {invoice.selected === category ? 'border border-primary-500' : ''} "
                                         on:click={() => {invoice.selected = category}}
                                    >
                                        <CategoryIcon size={48} scale={false} category={getCategoryByName(categories, category)} />
                                        <div class="flex flex-col">
                                            <span>{Money.format(values.benefit)}</span>
                                            <span class="text-gray-400">{Money.format(values.others)}</span>
                                        </div>
                                    </div>
                                {/if}
                            {/each}
                        {:else}
                            <div>Invoice has not been processed yet!</div>
                        {/if}
                    </div>
                </div>
            </TableBodyCell>
        </TableBodyRow>
    {/if}
{/if}