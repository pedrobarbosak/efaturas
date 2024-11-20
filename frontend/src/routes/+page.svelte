<script lang="ts">
    import { onMount } from "svelte";

    import {Alert, Avatar, Card, Checkbox, Input, Label} from 'flowbite-svelte';
    import { Button } from "flowbite-svelte";
    import { InfoCircleSolid } from 'flowbite-svelte-icons';

    import { Form } from "@/models/ui";
    import {translate} from "@/translations";
    import {goto} from "$app/navigation";
    import {formKey, user} from "@/store";

    let form: Form = new Form();

    onMount(() => {
        loadLocalStorage();
    });

    async function onFormSubmit() {
        form.isLoading = true
        const token = await form.onSubmit()
        form = form

        if (token === "") {
            return
        }

        clearLocalStorage();
        if (form.remember) {
            saveLocalStorage();
        }

        user.set({username: form.username.value, password: form.username.value, isLoggedIn: true, token: token})
        await goto("/invoices")
    }

    function loadLocalStorage() {
        const data = localStorage.getItem(formKey)
        if (data) {
            const { uname, pword, remember } = JSON.parse(data)
            if (remember) {
                form.username.value = uname;
                form.password.value = pword;
                form.remember = remember
            }
        }
    }

    function saveLocalStorage() {
        localStorage.setItem(formKey, JSON.stringify({uname: form.username.value, pword: form.password.value, remember: form.remember}))
    }


    function clearLocalStorage() {
        localStorage.removeItem(formKey)
    }

</script>

<div class="min-h-screen flex items-center justify-center">
    <Card class="w-full max-w-md p-8">
        <form class="space-y-6" on:submit={async() => {await onFormSubmit()}}>
            <div class="flex justify-center">
                <Avatar size="xl" src="/efaturas.webp" />
            </div>

            <div>
                <Label for="username" class="mb-2">{translate("label_nif")}</Label>
                <Input id="username" type="text" placeholder="111 222 333" bind:value={form.username.value} disabled={form.isLoading} required/>
            </div>

            <div>
                <Label for="password" class="mb-2">{translate("label_password")}</Label>
                <Input id="password" type="password" placeholder="••••••••" bind:value={form.password.value} disabled={form.isLoading} required/>
            </div>

            <div class="flex items-center justify-between">
                <Checkbox bind:checked={form.remember}>{translate("label_remember")}</Checkbox>
            </div>

            <Button type="submit" class="w-full">Sign in</Button>

            {#if form.error}
                <Alert color="red" border dismissable>
                    <InfoCircleSolid slot="icon" class="w-5 h-5" />
                    {form.error}
                </Alert>
            {/if}
        </form>
    </Card>
</div>
