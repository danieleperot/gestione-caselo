<script setup lang="ts">
import { useAuth } from "../composables/useAuth";

const { user, login, logout } = useAuth();

const handleLogin = () => {
    login("test@example.com", "Password123!")
        .then(() => console.log("Logged in successfully"))
        .catch((error) => console.error("Login failed", error));
};
</script>

<template>
    <div
        class="px-8 py-6 rounded-lg border-yellow-400 border-2 flex items-center justify-center"
    >
        <div v-if="!user" class="logged-out">
            <button
                @click="handleLogin"
                class="border rounded-md shadow px-6 py-2 hover:shadow-lg transition-all active:bg-gray-200 cursor-pointer"
            >
                Login as Test User
            </button>
        </div>
        <div v-else class="logged-in">
            <span class="user-info">Logged in as: {{ user.email }}</span>
            <button
                @click="logout"
                class="border rounded-md shadow px-6 py-2 hover:shadow-lg transition-all active:bg-gray-200 cursor-pointer"
            >
                Logout
            </button>
        </div>
    </div>
</template>
