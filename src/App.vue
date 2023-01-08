<template>
  <router-view />
</template>

<script lang="ts" setup>
import { useTheme } from "vuetify/lib/framework.mjs";

const theme = useTheme();

theme.global.name.value = "dark";

//
</script>
