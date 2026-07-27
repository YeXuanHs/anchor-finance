<template>
  <div />
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: [
    function () {
      return navigateTo('/admin/dashboard', { redirectCode: 301 })
    },
  ],
})
</script>
