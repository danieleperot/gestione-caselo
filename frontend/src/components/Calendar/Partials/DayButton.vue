<script setup lang="ts">
interface Props {
    busy?: boolean;
    day?: number | string;
    disabled?: boolean;
    outOfMonth?: boolean;
    selected?: boolean;
    today?: boolean;
}

withDefaults(defineProps<Props>(), {
    busy: false,
    day: 0,
    disabled: false,
    outOfMonth: false,
    selected: false,
    today: false,
});
</script>

<template>
    <div class="day-button w-10 h-10 shrink-0 flex items-center justify-center">
        <button
            :disabled="disabled || selected"
            class="day-button__inner relative rounded-full w-full h-full transition duration-150"
            :class="{
                'before:absolute before:bg-slate-500 active:bg-slate-500 before:rounded-full tex':
                    busy && !disabled,
                'bg-slate-50': !outOfMonth && !disabled,
                'text-slate-500': outOfMonth && !disabled,
                'text-slate-200 before:none': disabled,
                'hover:bg-slate-300 hover:text-slate-900':
                    !disabled && !selected,
                'border-2 border-slate-500': today,
                'bg-teal-100': selected,
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
