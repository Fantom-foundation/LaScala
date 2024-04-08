<template>
  <v-snackbar v-model="getSnackbarVisible" :timeout="-1" rounded location="bottom right" color="blue-lighten-5">
    <div v-for="(job, index) in getComparisonJobs" :key="index" class="d-flex align-center justify-space-between pb-2">
      <p class="comparison-job">
        {{ job.split(".")[0].split("_").join(" ") }}
      </p>
      <v-btn density="compact" variant="text" color="red" icon="mdi-trash-can-outline"
        @click="removeComparisonJob(job)"></v-btn>
    </div>

    <div class="snackbar-btn-wrapper">
      <v-btn color="indigo-accent-4" variant="flat" @click="sendComparisonJobs" :disabled="getComparisonJobs.length < 2">
        Send
      </v-btn>
      <v-btn color="red" variant="flat" @click="clearComparisonJobs">
        Clear
      </v-btn>
    </div>
  </v-snackbar>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
export default {
  computed: {
    ...mapGetters(["getComparisonJobs", "getSnackbarVisible"]),
  },

  methods: {
    ...mapActions(['clearComparisonJobs', 'removeComparisonJob', 'sendComparisonJobs'])
  }
}
</script>