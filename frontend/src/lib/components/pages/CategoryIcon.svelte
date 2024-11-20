<script lang="ts">
  import type { Category } from "@/models/models";
  import { translate } from "@/translations";
  import {Money} from "@/models/models";

  export let category: Category;
  export let size: number = 64;
  export let disabled: boolean = false;
  export let scale: boolean = true;
  export let showTotal: boolean = false
</script>


{#if category.total == null || !showTotal}

  <div class="flex p-2 cursor-pointer {disabled ? 'disabled' : '' }" title="{translate(category.name)}">
    <div class="icons icons-{size} {scale ? 'hover:scale-[1.25]' : `hover:scale-[$scale]` }" style="color: {category.color}">
      {String.fromCharCode(parseInt(category.unicode, 16))}
    </div>
  </div>

{:else}

  <div class="flex flex-col p-2 cursor-pointer items-center justify-items-center {disabled ? 'disabled' : '' } {scale ? 'hover:scale-[1.25]' : `hover:scale-[$scale]` }" title="{translate(category.name)}">
    <div class="icons icons-{size}" style="color: {category.color}">
      {String.fromCharCode(parseInt(category.unicode, 16))}
    </div>
    <span class="text-sm dark:text-white">
      {Money.format(category.total)} €
    </span>
  </div>

{/if}


<style>
  .disabled {
    filter: grayscale(100%);
    opacity: 0.5;
  }
</style>