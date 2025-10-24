<script setup lang="ts">
import { computed, ref } from "vue";
import MonthCalendar from "./components/Calendar/MonthCalendar.vue";
import HelloWorld from "./components/HelloWorld.vue";

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

const selected = ref<Date>(new Date());
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
    <main class="flex items-center justify-center grow shrink min-h-0">
        <div>
            <HelloWorld />

            <MonthCalendar
                v-model="selected"
                :minimum-date="minimumDateBeforeEvent"
                @view-changed="handleViewChanged"
            />

            <div class="text-xs p-4 text-center">
                <div>Selected: {{ selected.toISOString() }}</div>

                <div>
                    <strong>MONTH:</strong>
                    {{ viewMonth }}
                </div>

                <div>
                    <strong>RANGE:</strong>
                    {{ viewRange }}
                </div>
            </div>
        </div>
    </main>
</template>
