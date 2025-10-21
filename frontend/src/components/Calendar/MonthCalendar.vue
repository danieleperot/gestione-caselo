<script setup>
import DayButton from "./Partials/DayButton.vue";
import WeekDay from "./Partials/WeekDay.vue";
</script>

<template>
    <div class="border border-teal-200 p-4 rounded-2xl">
        <div class="flex items-center justify-between mb-6">
            <button class="w-10 h-10 flex items-center justify-center">
                &lt;
            </button>
            <div>Ottobre 2025</div>
            <button class="w-10 h-10 flex items-center justify-center">
                &gt;
            </button>
        </div>
        <div class="grid grid-cols-7 text-center gap-x-3 gap-y-2 text-sm">
            <WeekDay letter="L" full-day="Lunedí" />
            <WeekDay letter="M" full-day="Martedí" />
            <WeekDay letter="M" full-day="Mercoledí" />
            <WeekDay letter="G" full-day="Giovedí" />
            <WeekDay letter="V" full-day="Venerdí" />
            <WeekDay letter="S" full-day="Sabato" />
            <WeekDay letter="D" full-day="Domenica" />

            <DayButton day="28" out-of-month />
            <DayButton day="29" out-of-month />
            <DayButton day="30" out-of-month />
            <DayButton
                v-for="(_, key) in new Array(31)"
                :key="key"
                :day="key + 1"
                :busy="(key + 1) % 2 === 0"
            />
            <DayButton day="1" out-of-month />
        </div>
    </div>
</template>
