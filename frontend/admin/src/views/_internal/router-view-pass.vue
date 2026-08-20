<template>
  <router-view />
</template>

<script setup lang="ts">
/**
 * 透传组件 - 用于多级菜单的中间层级
 * 仅渲染 <router-view>，不添加任何布局或样式
 */
defineOptions({ name: 'RouterViewPass' })
</script>
