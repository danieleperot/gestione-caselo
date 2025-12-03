<script setup lang="ts">
import { computed, ref } from "vue";
import BookingForm from "@/components/Calendar/BookingForm.vue";
import MonthCalendar from "@/components/Calendar/MonthCalendar.vue";
import SelectedDate from "@/components/Calendar/SelectedDate.vue";
import { useGraphQL } from "@/composables/useGraphQL";
import { gql } from "graphql-request";

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
const isSubmitting = ref(false);
const submitError = ref<string | null>(null);
const submitSuccess = ref(false);

const { client } = useGraphQL();

const minimumDateBeforeEvent = computed<Date>(() => {
    const minimum = new Date();
    minimum.setUTCDate(minimum.getUTCDate() + 3);

    return minimum;
});

const handleViewChanged = ({ month, range }: ViewChangedPayload): void => {
    viewMonth.value = `Beginning: ${month.begin.toISOString()} | End: ${month.end.toISOString()}`;
    viewRange.value = `Beginning: ${range.begin.toISOString()} | End: ${range.end.toISOString()}`;
};

const SUBMIT_EVENT_BOOKING = gql`
    mutation SubmitEventBooking($input: EventBookingInput!) {
        submitEventBooking(input: $input) {
            success
            message
        }
    }
`;

const handleSubmit = async (event: Event): Promise<void> => {
    event.preventDefault();

    submitError.value = null;
    submitSuccess.value = false;

    const formData = new FormData(event.target as HTMLFormElement);

    const input = {
        fullName: formData.get("full_name") as string,
        association: formData.get("association") as string || undefined,
        email: formData.get("email") as string,
        phone: formData.get("phone") as string,
        description: formData.get("description") as string,
        date: formData.get("date") as string,
        acceptData: formData.get("accept_data") === "on",
    };

    isSubmitting.value = true;

    try {
        const response = await client.request<{
            submitEventBooking: { success: boolean; message: string };
        }>(SUBMIT_EVENT_BOOKING, { input });

        if (response.submitEventBooking.success) {
            submitSuccess.value = true;
            (event.target as HTMLFormElement).reset();
            selected.value = null;
        } else {
            submitError.value = response.submitEventBooking.message;
        }
    } catch (error) {
        console.error("Error submitting booking:", error);
        submitError.value = "Si è verificato un errore durante l'invio della richiesta. Riprova più tardi.";
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <div class="flex items-center bg-gray-100 py-6 text-gray-900 grow h-full">
        <form
            @submit="handleSubmit"
            class="w-full max-w-7xl mx-auto bg-white rounded-2xl shadow border-2 border-slate-300 py-12 px-8"
        >
            <div class="mb-12">
                <h2 class="font-bold text-2xl mb-3">Prenota il tuo evento</h2>
                <div class="text-slate-500">
                    Compila il modulo per richiedere la prenitazione. Una volta
                    inviato, ti manderemo una email di riepilogo. La
                    prenotazione
                    <strong>non è da considerarsi confermata</strong> fino a
                    quando non sarai ricontattato da un membro del gruppo.
                </div>
            </div>

            <div
                v-if="submitSuccess"
                class="mb-8 p-4 bg-green-100 border border-green-400 text-green-700 rounded-lg"
            >
                Richiesta inviata con successo! Riceverai una email di conferma a breve.
            </div>

            <div
                v-if="submitError"
                class="mb-8 p-4 bg-red-100 border border-red-400 text-red-700 rounded-lg"
            >
                {{ submitError }}
            </div>

            <div
                class="grid gap-8 items-start"
                :class="selected ? 'grid-cols-3' : 'grid-cols-2'"
            >
                <BookingForm />
                <div class="mt-6">
                    <input
                        type="date"
                        :value="selected?.toISOString().split('T')[0]"
                        name="date"
                        required
                        class="sr-only"
                    />
                    <MonthCalendar
                        v-model="selected"
                        :minimum-date="minimumDateBeforeEvent"
                        @view-changed="handleViewChanged"
                    />
                </div>
                <SelectedDate
                    v-if="selected"
                    :selected="selected"
                    class="mt-6"
                    @cancel="selected = null"
                />
            </div>
            <div class="flex items-center justify-center mt-12">
                <button
                    type="submit"
                    :disabled="isSubmitting"
                    class="rounded-lg text-xl bg-purple-500 text-white font-bold px-12 py-4 shadow hover:bg-purple-700 transition cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {{ isSubmitting ? "Invio in corso..." : "Invia richiesta di prenotazione" }}
                </button>
            </div>
        </form>
    </div>
</template>
