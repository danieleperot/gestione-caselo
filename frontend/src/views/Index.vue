<script setup lang="ts">
import { computed, ref } from "vue";
import MonthCalendar from "@/components/Calendar/MonthCalendar.vue";
import SelectedDate from "@/components/Calendar/SelectedDate.vue";

interface ViewChangedPayload {
    month: {
        begin: Date;
        end: Date;
    };
    range: {
        begin: Date;
        end: Date;
    };
}

const selected = ref<Date | null>(null);
const viewMonth = ref<string>("");
const viewRange = ref<string>("");

const minimumDateBeforeEvent = computed<Date>(() => {
    const minimum = new Date();
    minimum.setUTCDate(minimum.getUTCDate() + 3);

    return minimum;
});

const handleViewChanged = ({ month, range }: ViewChangedPayload): void => {
    viewMonth.value = `Beginning: ${month.begin.toISOString()} | End: ${month.end.toISOString()}`;
    viewRange.value = `Beginning: ${range.begin.toISOString()} | End: ${range.end.toISOString()}`;
};
</script>

<template>
    <header />
    <main class="md:flex md:items-center h-full py-4">
        <div
            class="md:flex content-center md:justify-center grow shrink min-h-0 md:w-full gap-12 md:p-6"
        >
            <MonthCalendar
                v-model="selected"
                :minimum-date="minimumDateBeforeEvent"
                class="grow md:max-w-xl justify-self-center w-full md:w-auto"
                @view-changed="handleViewChanged"
            />
            <SelectedDate
                v-if="selected"
                :selected="selected"
                class="grow md:max-w-xl fixed md:static bottom-0 left-0 right-0 bg-white w-full md:w-auto max-h-[80vh] overflow-y-auto"
                @cancel="selected = null"
            />
        </div>
    </main>
    <footer />
</template>
