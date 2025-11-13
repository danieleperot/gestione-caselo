<script setup lang="ts">
import { withDefaults } from "vue";

interface Props {
    label?: string;
    name: string;
    placeholder?: string;
    required?: boolean;
}

withDefaults(defineProps<Props>(), {
    label: "",
    placeholder: "",
    required: false,
});
</script>

<template>
    <label class="block">
        <div class="mb-1 ml-1">
            <strong>
                <span>{{ label }}</span>
            </strong>
            <span v-if="required" class="ml-1 text-slate-500">*</span>
        </div>
        <textarea
            :name="name"
            :required="required"
            :placeholder="placeholder"
            class="block border-2 border-slate-300 rounded-lg w-full py-2 px-4 placeholder-slate-400 text-slate-800"
        />
    </label>
</template>
