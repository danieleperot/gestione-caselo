<script setup>
defineProps({
    letter: {
        type: String,
        default: "",
    },
    fullDay: {
        type: String,
        default: "",
    },
});
</script>

<template>
    <div class="text-slate-500 mb-4" :arial-label="fullDay">{{ letter }}</div>
</template>
