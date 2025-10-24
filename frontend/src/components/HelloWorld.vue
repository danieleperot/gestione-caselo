<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useGraphQL } from "../composables/useGraphQL";
import { gql } from "graphql-request";

const { client } = useGraphQL();
const message = ref<string>("");
const loading = ref<boolean>(true);
const error = ref<string | null>(null);

const HELLO_QUERY = gql`
    query {
        hello {
            message
        }
    }
`;

onMounted(() => {
    client
        .request<{ hello: { message: string } }>(HELLO_QUERY)
        .then((data) => (message.value = data.hello.message))
        .catch((e) => {
            error.value = e instanceof Error ? e.message : "An error occurred";
        })
        .finally(() => (loading.value = false));
});
</script>

<template>
    <div class="p-4">
        <div v-if="loading">Loading...</div>
        <div v-else-if="error" class="text-red-500">Error: {{ error }}</div>
        <div v-else class="text-2xl font-bold">{{ message }}</div>
    </div>
</template>
