<script setup lang="ts">
import { computed, ref, watch } from "vue";
import XMark from "../Icons/XMark.vue";

const props = defineProps<{
    selected: Date;
}>();

defineEmits<{ cancel: [] }>();

type SlotType = "morning" | "afternoon" | "evening" | "full_day";
const selectedSlot = ref<SlotType | null>(null);

const day = computed<string>(() => {
    return `${props.selected.getUTCDate()}`;
});

const month = computed<string>(() => {
    return props.selected.toLocaleString("it-IT", { month: "long" });
});

const year = computed<string>(() => {
    return `${props.selected.getUTCFullYear()}`;
});

const slots = {
    morning: {
        label: "Mattina",
        hint: "Indicativamente dalle 06:30 alle 13:30",
    },
    afternoon: {
        label: "Pomeriggio",
        hint: "Indicativamente dalle 14:00 alle 20:00",
    },
    evening: {
        label: "Sera",
        hint: "Indicativamente dalle 20:30 alle 00:00",
    },
    full_day: {
        label: "Giornata intera",
        hint: "Indicativamente dalle 06:30 alle 00:00",
    },
};

const selectedIsoDate = computed(() => props.selected.toISOString());

watch(selectedIsoDate, () => {
    selectedSlot.value = null;
});
</script>

<template>
    <div class="rounded-t-2xl md:rounded-b-2xl border-2 border-slate-300">
        <div class="flex items-center justify-end">
            <button
                aria-label="Scegli un'altra data"
                class="w-12 h-12 flex items-center justify-center rounded-full cursor-pointer"
                @click="$emit('cancel')"
            >
                <div
                    class="w-10 h-10 hover:shadow-md transition-all duration-300 hover:bg-slate-50/80 rounded-full flex items-center justify-center bg-white relative"
                >
                    <XMark class="size-4" />
                </div>
            </button>
        </div>
        <div class="pt-2 pb-2 md:pb-6 px-2 md:px-6 space-y-6">
            <div class="text-center text-lg">
                <strong>{{ day }} {{ month }} {{ year }}</strong>
            </div>
            <div class="grid gap-2">
                <button
                    v-for="({ label, hint }, key) in slots"
                    :key="key"
                    class="border-2 border-slate-200 grid gap-0 rounded-lg w-full p-2 hover:shadow-md transition-all cursor-pointer active:bg-purple-50"
                    :class="{
                        'bg-purple-200': selectedSlot === key,
                        'bg-white': selectedSlot !== key,
                    }"
                    :disabled="selectedSlot === key"
                    @click="selectedSlot = key"
                >
                    <div>{{ label }}</div>
                    <div class="text-xs">
                        {{ hint }}
                    </div>
                </button>
            </div>
        </div>
    </div>
</template>
