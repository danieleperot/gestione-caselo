<script setup lang="ts">
interface Props {
    letter?: string;
    fullDay?: string;
}

withDefaults(defineProps<Props>(), {
    letter: "",
    fullDay: "",
});
</script>

<template>
    <div class="text-slate-500 mb-4" :arial-label="fullDay">{{ letter }}</div>
</template>
