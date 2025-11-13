<script setup lang="ts">
import { computed, ref, onMounted } from "vue";

import DayButton from "./Partials/DayButton.vue";
import WeekDay from "./Partials/WeekDay.vue";

import ChevronLeft from "../Icons/ChevronLeft.vue";
import ChevronRight from "../Icons/ChevronRight.vue";

interface DateInfo {
    day: number;
    date: Date;
}

const props = withDefaults(defineProps<{ minimumDate?: Date }>(), {
    minimumDate: () => new Date(),
});

const emit = defineEmits<{
    viewChanged: [
        payload: {
            month: {
                begin: Date;
                end: Date;
            };
            range: {
                begin: Date;
                end: Date;
            };
        },
    ];
}>();

const selectedDate = defineModel<Date | null>({ default: () => null });
const viewingDate = ref<Date>(new Date(selectedDate.value || new Date()));

const selectedMidnight = computed<Date>(() => {
    const midnight = new Date(selectedDate.value || new Date());
    midnight.setUTCHours(0, 0, 0, 0);

    return midnight;
});

const minimumDateBeforeEvent = computed<Date>(() => {
    const minimum = new Date(props.minimumDate);
    minimum.setUTCHours(0, 0, 0, 0);

    return minimum;
});

const viewingMonth = computed<number>(() => {
    return viewingDate.value.getUTCMonth();
});

const viewingYear = computed<number>(() => {
    return viewingDate.value.getUTCFullYear();
});

const viewingMonthLabel = computed<string>(() => {
    return viewingDate.value.toLocaleString("it-IT", { month: "long" });
});

const previousMonth = computed<Date>(() => {
    const previousMonth = new Date(viewingDate.value);
    previousMonth.setUTCMonth(viewingMonth.value - 1);
    previousMonth.setUTCDate(1);

    return previousMonth;
});

const firstWeekDayOfTheMonth = computed<number>(() => {
    const firstDay = new Date(viewingDate.value);
    firstDay.setUTCDate(1);

    return firstDay.getUTCDay();
});

const canGoBack = computed<boolean>(() => {
    return (
        previousMonth.value.getUTCFullYear() >
            minimumDateBeforeEvent.value.getUTCFullYear() ||
        previousMonth.value.getUTCMonth() >=
            minimumDateBeforeEvent.value.getUTCMonth()
    );
});

const daysInCurrentMonth = computed<number>(() => {
    // Very neat trick suggested by AI. By setting the day to 0,
    // JS actually rolls back to the last day of the previous month
    const monthDate = new Date(viewingDate.value);
    monthDate.setUTCMonth(monthDate.getUTCMonth() + 1);
    monthDate.setUTCDate(0);

    return monthDate.getUTCDate();
});

const previousDaysToRender = computed<DateInfo[]>(() => {
    const size =
        firstWeekDayOfTheMonth.value > 0 ? firstWeekDayOfTheMonth.value - 1 : 6;

    return new Array(size)
        .fill({})
        .map((_, index) => {
            const date = new Date(viewingDate.value);
            date.setUTCMonth(date.getUTCMonth());
            date.setUTCDate(0 - index);
            date.setUTCHours(0, 0, 0, 0);

            return { day: date.getDate(), date };
        })
        .reverse();
});

const nextDaysToRender = computed<DateInfo[]>(() => {
    const maxDays = 6 /* weeks */ * 7; /* days a week */
    const renderedUntilNow =
        previousDaysToRender.value.length + daysInCurrentMonth.value;
    const size = maxDays - renderedUntilNow;

    if (isNaN(size)) {
        return new Array<DateInfo>();
    }

    return new Array(maxDays - renderedUntilNow).fill({}).map((_, index) => {
        const date = new Date(viewingDate.value);
        date.setUTCMonth(date.getUTCMonth() + 1);
        date.setUTCDate(index + 1);
        date.setUTCHours(0, 0, 0, 0);

        return { day: index + 1, date };
    });
});

const prevMonth = (date: Date | null): void => changeMonth(-1, date);
const nextMonth = (date: Date | null): void => changeMonth(+1, date);

const changeMonth = (toAdd: number, date: Date | null): void => {
    if (date) {
        selectedDate.value = date;
    }

    const newDate = new Date(viewingDate.value);
    newDate.setUTCMonth(viewingMonth.value + toAdd);

    viewingDate.value = newDate;

    emitViewChanged();
};

const emitViewChanged = (): void => {
    const beginMonth = new Date(viewingDate.value);
    beginMonth.setUTCDate(1);
    beginMonth.setUTCHours(0, 0, 0, 0);

    const endMonth = new Date(viewingDate.value);
    endMonth.setUTCMonth(endMonth.getUTCMonth() + 1);
    endMonth.setUTCDate(0);
    endMonth.setUTCHours(23, 59, 59, 999);

    let beginRange = new Date(beginMonth);
    if (previousDaysToRender.value.length) {
        beginRange = new Date(previousDaysToRender.value[0].date);
        beginRange.setUTCHours(0, 0, 0, 0);
    }

    let endRange = new Date(endMonth);
    if (nextDaysToRender.value.length) {
        const lastDay = nextDaysToRender.value.at(-1);
        if (lastDay) {
            endRange = new Date(lastDay.date);
            endRange.setUTCHours(23, 59, 59, 999);
        }
    }

    emit("viewChanged", {
        month: { begin: beginMonth, end: endMonth },
        range: { begin: beginRange, end: endRange },
    });
};

const datesToDisplay = computed<DateInfo[]>(() => {
    return new Array(daysInCurrentMonth.value).fill({}).map((_, index) => {
        const date = new Date(viewingDate.value);
        date.setUTCDate(index + 1);
        date.setUTCHours(0, 0, 0, 0);

        return { day: index + 1, date };
    });
});

const now = ref<Date>(new Date());
now.value.setUTCHours(0, 0, 0, 0);

onMounted(() => {
    const lastOfMonth = new Date(viewingDate.value);
    lastOfMonth.setUTCMonth(lastOfMonth.getUTCMonth() + 1);
    lastOfMonth.setUTCDate(0); // Roll back to current month on last day
    lastOfMonth.setUTCHours(23, 59, 59, 999);

    const minimum = new Date(minimumDateBeforeEvent.value);
    minimum.setUTCHours(0, 0, 0, 0);

    // On purpose less or equal
    if (lastOfMonth.getTime() <= minimum.getTime()) {
        viewingDate.value = new Date(minimumDateBeforeEvent.value);
    }

    // On purpose less than, not equal
    if (
        selectedDate.value &&
        selectedMidnight.value.getTime() < minimum.getTime()
    ) {
        selectedDate.value = new Date(minimumDateBeforeEvent.value);
    }

    emitViewChanged();
});

const isBusy = (date: Date): boolean => {
    return date.getDate() % 3 === 0;
};
</script>

<template>
    <div class="border-2 border-slate-300 p-4 rounded-2xl text-slate-900">
        <div class="flex items-center justify-between mb-6">
            <button
                type="button"
                :disabled="!canGoBack"
                class="w-10 h-10 flex items-center justify-center disabled:text-slate-400 rounded-full bg-transparent cursor-pointer"
                :class="{
                    'hover:bg-slate-50 active:bg-slate-50 transition':
                        canGoBack,
                }"
                @click="prevMonth(null)"
            >
                <ChevronLeft aria-label="Mese precedente" />
            </button>
            <div class="capitalize">
                {{ viewingMonthLabel }} {{ viewingYear }}
            </div>
            <button
                type="button"
                class="w-10 h-10 flex items-center justify-center disabled:text-slate-400 rounded-full bg-transparent hover:bg-slate-50 active:bg-slate-50 transition cursor-pointer"
                @click="nextMonth(null)"
            >
                <ChevronRight aria-label="Mese successivo" />
            </button>
        </div>
        <div
            class="grid grid-cols-7 text-center gap-x-1 gap-y-2 text-sm sm:gap-x-3 justify-items-center"
        >
            <WeekDay letter="L" full-day="Lunedí" />
            <WeekDay letter="M" full-day="Martedí" />
            <WeekDay letter="M" full-day="Mercoledí" />
            <WeekDay letter="G" full-day="Giovedí" />
            <WeekDay letter="V" full-day="Venerdí" />
            <WeekDay letter="S" full-day="Sabato" />
            <WeekDay letter="D" full-day="Domenica" />

            <DayButton
                v-for="{ date, day } in previousDaysToRender"
                :key="day"
                :busy="isBusy(date)"
                :day="day"
                :disabled="
                    isBusy(date) ||
                    !canGoBack ||
                    date.getTime() < minimumDateBeforeEvent.getTime()
                "
                out-of-month
                @click="prevMonth(date)"
            />

            <DayButton
                v-for="{ date, day } in datesToDisplay"
                :key="day"
                :busy="isBusy(date)"
                :day="day"
                :disabled="
                    isBusy(date) ||
                    date.getTime() < minimumDateBeforeEvent.getTime()
                "
                :selected="
                    Boolean(
                        selectedDate &&
                            date.getTime() === selectedMidnight.getTime(),
                    )
                "
                :today="date.getTime() === now.getTime()"
                @click="selectedDate = new Date(date)"
            />

            <DayButton
                v-for="{ date, day } in nextDaysToRender"
                :key="day"
                :busy="isBusy(date)"
                :disabled="isBusy(date)"
                :day="day"
                out-of-month
                @click="nextMonth(date)"
            />
        </div>
    </div>
</template>
