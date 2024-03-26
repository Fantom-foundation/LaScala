import { createStore } from 'vuex'
import axios from 'axios'

// Create a new store instance.
const store = createStore({
  state: {
    runGroups: null,
    runGroupJobs: {},
    jobReport: ""
  },
  getters: {
    getAllRunGroups: state => state.runGroups,
    getRunGroupJobs: state => runGroup => state.runGroupJobs[runGroup],
    getJobReport: state => state.jobReport,
  },
  mutations: {
    addRunGroups(state, runs) {
      state.runGroups = runs;
      runs.forEach(run => { state.runGroupJobs[run] = null })
    },
    addRunGroupJobs(state, payload) {
      state.runGroupJobs[payload.group] = payload.jobs;
    },
    updateJobReport(state, report) {
      state.jobReport = report
    }
  },
  actions: {
    async fetchAllRunGroups(context) {
      try {
        let response = await axios.get(`${import.meta.env.VITE_API_URL}/static/`)
        let runGroups = response.data.split('href').map(substring => {
          let result = substring.match(/[A-Za-z]+_[0-9]+/);

          if (!result) {
            return
          }

          return result[0]
        }).filter(val => !!val)

        context.commit('addRunGroups', runGroups)
      } catch (error) {
        console.log(error);
      }
    },

    async fetchAllRunGroupJobs(context, runGroupName) {
      if (!runGroupName) {
        return
      }

      try {
        let response = await axios.get(`${import.meta.env.VITE_API_URL}/master/${runGroupName}/`);
        let jobs = response.data.split('href').map(substring => {
          let result = substring.match(/([A-Za-z0-9]+(_[A-Za-z0-9]+)+)\.[A-Za-z]+/);

          if (!result) {
            return
          }

          return result[0]
        }).filter(val => !!val)

        context.commit('addRunGroupJobs', { jobs: jobs, group: runGroupName })
      } catch (error) {
        console.log(error);
      }
    },

    async fetchReport(context, payload) {
      const { runGroupName, jobName } = payload;

      if (!runGroupName || !jobName) {
        return
      }

      try {
        context.commit('updateJobReport', "");

        let response = await axios.get(`${import.meta.env.VITE_API_URL}/report/${runGroupName}/${jobName}`);

        context.commit('updateJobReport', response.data)
      } catch (error) {
        console.log(error);
      }
    }
  }
})

export default store