	
<script lang="ts">
    import {Card} from 'flowbite-svelte';
    import {getInvoiceSSE} from "@/api";
    import {Category, Invoice} from "@/models/models";
    import {onDestroy, onMount} from "svelte";
    import {decodeHTML} from "@/utils";

    let eventSource : EventSource;
    export let invoice : Invoice = new Invoice();
    export let categories : Category[] = [];

    onMount(()=> {
      eventSource = getInvoiceSSE(invoice.id)
      eventSource.onmessage = (event) => {
          invoice = new Invoice(JSON.parse(event.data));
      }
    })

    onDestroy(() => {
      if(eventSource) {
        eventSource.close()
      }
    })

    function getCategoryByName(category: string) : Category {
      for (let i = 0; i < categories.length; i++) {
        const element = categories[i];
        if (element.name == category) {
          return element;
        }
      }

      return new Category();
    }

  </script>
   
<Card class="max-w-sm rounded overflow-hidden shadow-lg m-2 hover:scale-[1.1] cursor-pointer">
    <div>
      <div class="flex items-center">
        <div size={32} category={getCategoryByName(invoice.activity.category)}/> {decodeHTML(invoice?.document?.number)}
      </div>
      <div>
        <p>
          {invoice?.issuer.nif}
        </p>{decodeHTML(invoice?.issuer.name)}
      </div>
    </div>
    <div>
      <p>Card Content</p>
    </div>
    <div>
      <p>Card Footer</p>
    </div>
</Card>
<!--   
<div class="flex items-center justify-center border rounded p-2">
  <div>
    <CategoryIcon size={28} category={getCategoryByName(invoice.activity.category)}/>
  </div>
  <div>
    {decodeHTML(invoice?.document?.number)}
  </div>

  <div>
    {invoice?.issuer.nif} {decodeHTML(invoice?.issuer.name)}
  </div>
</div> -->