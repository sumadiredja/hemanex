<template>
  <v-app-bar :elevation="2" :collapse="globalStore.isCollapse()">
    <v-app-bar-title>Hemanex CLI</v-app-bar-title>

    <v-spacer></v-spacer>
    <template v-slot:append>
      <v-btn icon="mdi-menu" @click="globalStore.changeCollapse()"></v-btn>

      <v-btn icon="mdi-brightness-4" @click="changeTheme()"></v-btn>
    </template>
  </v-app-bar>
</template>

<script lang="ts" setup>
import { useTheme } from "vuetify/lib/framework.mjs";
import { useGlobalStore } from "../../store/app";
const globalStore = useGlobalStore();
const theme = useTheme();

function changeTheme() {
  theme.global.name.value = theme.global.current.value.dark ? "light" : "dark";
}
</script>
