<template>
  <v-tabs v-model="tab" align-tabs="left" color="indigo-accent-4">
    <v-tab v-for="(job, index) in jobs" :key="`tab_${index}`" :value="index + 1" @click="handleFetchReport(job)">
      {{ job.split(".")[0].split("_").join(" ") }}
    </v-tab>
  </v-tabs>
  <v-window v-model="tab">
    <template v-if="jobs">
      <v-window-item v-for="(job, index) in jobs" :key="`tab-window_${index}`" :value="index + 1">
        <v-container fluid>
          <v-row v-if="!!report">
            <v-col cols="12">
              <v-checkbox class="checkbox" label="Pick for comparison" :value="`${runGroupName}/${job}`"
                color="indigo-accent-4" v-model="comparison">
              </v-checkbox>
              <div v-html="report"></div>
            </v-col>
          </v-row>
          <v-row v-else>
            <v-col cols="12">
              <div class="centered">
                <v-progress-circular :size="70" :width="7" color="indigo-accent-4" indeterminate>
                </v-progress-circular>
              </div>
            </v-col>
          </v-row>
        </v-container>
      </v-window-item>
    </template>
  </v-window>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'

export default {
  props: {
    runGroupName: String
  },

  data() {
    return {
      tab: 1,
      comparison: []
    }
  },

  computed: {
    ...mapGetters(["getRunGroupJobs", "getJobReport", "getComparisonJobs"]),

    report() {
      return this.getJobReport;
    },

    jobs() {
      return this.getRunGroupJobs(this.runGroupName)
    },
  },

  watch: {
    runGroupName() {
      this.handleFetchAllRunGroupJobs()
    },

    comparison(val) {
      this.updateComparisonJobs(val)
    },

    getComparisonJobs(val) {
      this.comparison = val;
    }
  },

  methods: {
    ...mapActions(["fetchAllRunGroupJobs", "fetchReport", "updateComparisonJobs"]),

    async handleFetchAllRunGroupJobs() {
      await this.fetchAllRunGroupJobs(this.runGroupName);
      this.tab = 1;
      this.handleFetchReport(this.jobs[this.tab]);
    },

    async handleFetchReport(jobName) {
      if (!jobName) {
        return
      }

      await this.fetchReport({ runGroupName: this.runGroupName, jobName: jobName.split(".")[0] })
    }
  },

  mounted() {
    this.handleFetchAllRunGroupJobs()
  }
}
</script>

<style>
.centered {
  display: flex;
  justify-content: center;
  margin: 60px 0;
}
</style>
