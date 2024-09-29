<template>
  <div>
    <!-- Button to call API -->
    <button @click="newContainer">New</button>

    <!-- Display the response -->
    <div v-if="apiResponse">
      <h3>API Response:</h3>
      <pre>{{ apiResponse }}</pre>
    </div>

    <!-- Handle loading state -->
    <div v-if="loading">
      Loading...
    </div>

    <!-- Handle error -->
    <div v-if="error">
      <p>Error: {{ error }}</p>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

// Define the type for the API response
interface NewResponse {
  // Adjust this to match the structure of the API response
  name: string;
  created: string;
}
interface StatusResponse {
  name: string;
  ready: boolean;
}

export default defineComponent({
  data() {
    return {
      newResponse: null as NewResponse | null,  // Nullable type for API response
      loading: false,                           // Loading state
      error: null as string | null              // Nullable type for error message
    };
  },
  methods: {
    async newContainer() {
      this.loading = true;
      this.error = null;
      this.apiResponse = null;

      try {
        // You can type the API response if needed
        const response = await axios.get<ApiResponse>('/api/v0/ct/new');
        this.apiResponse = response.data;  // Store the typed API response
      } catch (err: any) {
        // Handle errors gracefully and extract error message
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }
    }
    async isReady() {
      this.loading = true;
      this.error = null;
      this.apiResponse = null;

      try {
        // You can type the API response if needed
        const response = await axios.get<ApiResponse>('/api/v0/ct/new');
        this.apiResponse = response.data;  // Store the typed API response
      } catch (err: any) {
        // Handle errors gracefully and extract error message
        this.error = err.response?.data?.message || err.message;
      } finally {
        this.loading = false;
      }     
    }
  }
});
</script>

