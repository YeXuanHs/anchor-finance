<template>
  <ElDialog
    v-model="visible"
    :title="dialogType === 'add' ? $t('roleEditDialog.addRole') : $t('roleEditDialog.editRole')"
    width="30%"
    align-center
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="120px">
      <ElFormItem :label="$t('roleEditDialog.roleName')" prop="roleName">
        <ElInput v-model="form.roleName" :placeholder="$t('roleEditDialog.inputRoleName')" />
      </ElFormItem>
      <ElFormItem :label="$t('roleEditDialog.roleCode')" prop="roleCode">
        <ElInput v-model="form.roleCode" :placeholder="$t('roleEditDialog.inputRoleCode')" />
      </ElFormItem>
      <ElFormItem :label="$t('roleEditDialog.description')" prop="description">
        <ElInput
          v-model="form.description"
          type="textarea"
          :rows="3"
          :placeholder="$t('roleEditDialog.inputDescription')"
        />
      </ElFormItem>
      <ElFormItem :label="$t('roleEditDialog.enabled')">
        <ElSwitch v-model="form.enabled" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="handleClose">{{ $t('common.cancel') }}</ElButton>
      <ElButton type="primary" @click="handleSubmit">{{ $t('common.submit') }}</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import type { FormInstance, FormRules } from 'element-plus'
  import request from '@/utils/http'
  import { $t } from '@/locales'

  type RoleListItem = Api.SystemManage.RoleListItem

  interface Props {
    modelValue: boolean
    dialogType: 'add' | 'edit'
    roleData?: RoleListItem
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    dialogType: 'add',
    roleData: undefined
  })

  const emit = defineEmits<Emits>()

  const formRef = ref<FormInstance>()

  /**
   * 弹窗显示状态双向绑定
   */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  /**
   * 表单验证规则
   */
  const rules = reactive<FormRules>({
    roleName: [
      { required: true, message: $t('roleEditDialog.inputRoleName'), trigger: 'blur' },
      { min: 2, max: 20, message: $t('roleEditDialog.nameLength'), trigger: 'blur' }
    ],
    roleCode: [
      { required: true, message: $t('roleEditDialog.inputRoleCode'), trigger: 'blur' },
      { min: 2, max: 50, message: $t('roleEditDialog.codeLength'), trigger: 'blur' }
    ],
    description: [{ required: true, message: $t('roleEditDialog.inputDescription'), trigger: 'blur' }]
  })

  /**
   * 表单数据
   */
  const form = reactive<RoleListItem>({
    roleId: 0,
    roleName: '',
    roleCode: '',
    description: '',
    createTime: '',
    enabled: true
  })

  /**
   * 监听弹窗打开，初始化表单数据
   */
  watch(
    () => props.modelValue,
    (newVal) => {
      if (newVal) initForm()
    }
  )

  /**
   * 监听角色数据变化，更新表单
   */
  watch(
    () => props.roleData,
    (newData) => {
      if (newData && props.modelValue) initForm()
    },
    { deep: true }
  )

  /**
   * 初始化表单数据
   * 根据弹窗类型填充表单或重置表单
   */
  const initForm = () => {
    if (props.dialogType === 'edit' && props.roleData) {
      Object.assign(form, props.roleData)
    } else {
      Object.assign(form, {
        roleId: 0,
        roleName: '',
        roleCode: '',
        description: '',
        createTime: '',
        enabled: true
      })
    }
  }

  /**
   * 关闭弹窗并重置表单
   */
  const handleClose = () => {
    visible.value = false
    formRef.value?.resetFields()
  }

  /**
   * 提交表单
   * 验证通过后调用接口保存数据
   */
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      if (props.dialogType === 'add') {
        await request.post({ url: '/api/admin/rbac/roles', data: { ...form } })
      } else {
        await request.put({ url: `/api/admin/rbac/roles/${form.roleId}`, data: { ...form } })
      }
      const message = props.dialogType === 'add' ? $t('roleEditDialog.addSuccess') : $t('roleEditDialog.editSuccess')
      ElMessage.success(message)
      emit('success')
      handleClose()
    } catch (error) {
      console.log($t('roleEditDialog.formValidateFailed') + ':', error)
    }
  }
</script>
