<script setup lang="ts">
import CheckMark from "../Icons/CheckMark.vue";
import TextField from "../Inputs/TextField.vue";
import TextareaField from "../Inputs/TextareaField.vue";
</script>

<template>
    <div class="space-y-5">
        <TextField
            name="full_name"
            label="Nome e cognome"
            placeholder="Mario Rossi"
            required
        />
        <TextField
            name="association"
            label="Associazione"
            placeholder="Associazione ricreativa culturale"
        />
        <TextField
            name="email"
            label="Indirizzo Email"
            placeholder="mario.rossi@example.com"
            required
            type="email"
        />
        <TextField
            name="phone"
            label="Numero di telefono"
            placeholder="0123 456789"
            required
            type="tel"
        />
        <TextareaField
            name="description"
            label="Descrizione dell'evento"
            required
            placeholder="Evento di 15 persone per un compleanno"
        />
        <label class="group flex items-center">
            <input
                name="accept_data"
                required
                type="checkbox"
                class="sr-only peer"
            />
            <div
                aria-hidden
                class="cursor-pointer w-10 h-10 border-2 border-slate-300 rounded-lg peer-checked:text-purple-400 text-white flex items-center justify-center"
            >
                <CheckMark class="size-8" />
            </div>
            <div class="pl-4">
                Acconsento al
                <a
                    href="/trattamento-dei-dati"
                    target="_blank"
                    class="cursor-pointer text-purple-600 trasition hover:text-purple-800"
                >
                    trattamento dei dati personali
                </a>
                (obbligatorio)
            </div>
        </label>
    </div>
</template>
