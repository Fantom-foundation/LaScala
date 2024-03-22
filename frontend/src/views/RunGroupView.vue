<template>
  <h2 class="mb-3 mt-0 h4">{{ heading }}</h2>

  <v-card class="pa-4 rounded-lg" rounded>
    <template v-if="renderJobList">
      <RunGroupJobs :runGroupName="runGroupName"></RunGroupJobs>
    </template>
  </v-card>
</template>

<script>
import { mapGetters } from 'vuex'
import RunGroupJobs from '@/components/RunGroupJobs.vue'

export default {
  props: {
    runGroupName: String
  },

  components: {
    RunGroupJobs
  },

  computed: {
    ...mapGetters(["getAllRunGroups"]),

    renderJobList() {
      return !!this.runGroupName
    },

    heading() {
      let capitalized = this.runGroupName.charAt(0).toUpperCase() + this.runGroupName.slice(1)
      return capitalized.split("_").join(" ")
    }
  }
}
</script>

<style></style>