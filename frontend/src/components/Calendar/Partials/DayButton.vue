<script setup>
defineProps({
    day: {
        type: Number,
        default: 0,
    },
    busy: {
        type: Boolean,
        default: false,
    },
    outOfMonth: {
        type: Boolean,
        default: false,
    },
});
</script>

<template>
    <div class="day-button w-10 h-10 shrink-0 flex items-center justify-center">
        <button
            class="day-button__inner relative rounded-full w-full h-full transition-all duration-300"
            :class="{
                'before:absolute before:bg-gray-500 active:bg-gray-500 before:rounded-full':
                    busy,
                'bg-gray-50 hover:bg-gray-300': !outOfMonth,
                'text-gray-300': outOfMonth,
            }"
        >
            {{ day }}
        </button>
    </div>
</template>

<style scoped>
.day-button__inner::before {
    height: 5px;
    left: 18px;
    top: 30px;
    width: 5px;
}
</style>
