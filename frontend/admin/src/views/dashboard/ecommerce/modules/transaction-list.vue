<template>
  <ArtDataListCard
    class="mb-5 max-sm:mb-4"
    :maxCount="4"
    :list="dataList"
    title="最近活动"
    subtitle="订单处理状态"
    :showMoreButton="true"
    @more="handleMore"
  />
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import request from '@/utils/http'

  interface TransactionItem {
    title: string
    status: string
    time: string
    class: string
    icon: string
  }

  const dataList = ref<TransactionItem[]>([])

  const fetchRecentActivity = async () => {
    try {
      const res = await request.get({ url: '/api/admin/dashboard/recent-activity' })
      if (Array.isArray(res)) {
        dataList.value = res
      } else {
        dataList.value = []
      }
    } catch {
      dataList.value = []
    }
  }

  const router = useRouter()

  const handleMore = (): void => {
    router.push('/finance/orders/list')
  }

  onMounted(() => {
    fetchRecentActivity()
  })
</script>
