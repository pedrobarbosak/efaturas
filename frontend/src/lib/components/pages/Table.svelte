<script lang="ts">
    import { translate } from '@/translations';
    import {type Category, Invoice, Money} from '@/models/models';
    import { Checkbox, Input, Table, TableBody, TableHead, TableHeadCell} from 'flowbite-svelte';
    import TableRow from './TableRow.svelte';
    import { CheckPlusCircleOutline, FileCirclePlusOutline, FileSearchOutline, SearchOutline } from "flowbite-svelte-icons";
    import { SpeedDial, SpeedDialButton } from 'flowbite-svelte';
    import {onMount} from "svelte";

    export let invoices: Invoice[] = [];
    export let categories: Category[] = [];

    let search: string = "";
    let expandAll: boolean = false;

</script>

<div class="flex p-4 justify-between">
  <div class="flex items-center p-2 w-full">
    <Input id="search" type="search" placeholder="search" bind:value={search}>
      <SearchOutline slot="left" class="w-5 h-5" />
    </Input>
  </div>

  <div class="flex items-center p-2 justify-items-end justify-end">
    <span class="p-2 min-w-max dark:text-white">
      {#if (expandAll)} {translate("collapse_all")} {:else} {translate("expand_all")} {/if}
    </span>
    <Checkbox bind:checked={expandAll}/>
  </div>
</div>

<Table hoverable={true} >
  <TableHead>
    <TableHeadCell>{translate("invoice_sector")}</TableHeadCell>
    <TableHeadCell>{translate("invoice_date")}</TableHeadCell>
    <TableHeadCell>{translate("invoice_number")}</TableHeadCell>
    <TableHeadCell>{translate("invoice_issuer")}</TableHeadCell>
    <TableHeadCell>{translate("invoice_possible_sectors")}</TableHeadCell>
    <TableHeadCell>{translate("invoice_value")}</TableHeadCell>
<!--    <TableHeadCell>{translate("invoice_benefit")}</TableHeadCell>-->
  </TableHead>
  <TableBody tableBodyClass="divide-y">
    {#each invoices as invoice}
      <TableRow bind:invoice={invoice} {categories} isOpen={expandAll} hide={!Invoice.search(invoice, search)}/>
    {/each}
  </TableBody>
</Table>

<SpeedDial >
  <SpeedDialButton name="{translate('btn_rescan')}" on:click={() => {console.log("A")}} >
    <FileSearchOutline class="w-6 h-6" />
  </SpeedDialButton>
  <SpeedDialButton name="{translate('btn_update')}" on:click={() => {console.log("B")}} >
    <FileCirclePlusOutline class="w-6 h-6" />
  </SpeedDialButton>
  <SpeedDialButton name="{translate('btn_apply_changes')}" on:click={() => {console.log("C")}} >
    <CheckPlusCircleOutline class="w-6 h-6" />
  </SpeedDialButton>
</SpeedDial>

