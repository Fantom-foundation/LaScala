<template>
  <v-app>
    <v-app-bar color="indigo-accent-4" prominent>
      <v-app-bar-nav-icon variant="text" @click.stop="drawer = !drawer">
      </v-app-bar-nav-icon>

      <v-toolbar-title>Amneris - LaScala reporting app</v-toolbar-title>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" location="left" temporary>
      <v-list>
        <v-list-item link @click="redirect({ name: 'home' })">
          Home
        </v-list-item>
        <v-list-item v-for="(runGroup, index) in runGroups" :key="index" link
          @click="redirect({ name: 'runGroup', params: { runGroupName: getAllRunGroups[index] } })">
          {{ runGroup }}
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main class="pa-6 pt-10 mt-16 bg-blue-lighten-5">
      <RouterView />
    </v-main>
  </v-app>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'

export default {
  data() {
    return {
      drawer: false
    }
  },

  computed: {
    ...mapGetters(["getAllRunGroups"]),

    runGroups() {
      if (!this.getAllRunGroups) {
        return []
      }

      return this.getAllRunGroups.map(group => {
        let capitalized = group.charAt(0).toUpperCase() + group.slice(1)
        return capitalized.split("_").join(" ")
      })
    }
  },

  methods: {
    ...mapActions(["fetchAllRunGroups"]),

    async handleFetchAllRunGroups() {
      await this.fetchAllRunGroups()
    },

    redirect(route) {
      this.$router.push(route);
    }
  },

  mounted() {
    this.fetchAllRunGroups()
  },
}
</script>

<style scoped></style>
