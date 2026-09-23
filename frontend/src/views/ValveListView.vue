<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../services/api'
import { FilterMatchMode } from '@primevue/core/api'
import FileUpload from 'primevue/fileupload'
import { uploadValveImage, uploadValveDatasheet } from '../services/api'


import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import ConfirmDialog from 'primevue/confirmdialog'
import { useConfirm } from 'primevue/useconfirm'
import Toast from 'primevue/toast'
import { useToast } from 'primevue/usetoast'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import { useRouter } from 'vue-router'
import { useValveStore } from '../stores/valve'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const valveStore = useValveStore()
const authStore = useAuthStore()

function viewDetail(valve) {
  valveStore.selectValve(valve)
  router.push({ name: 'valve-details', params: { id: valve.id } })
}

const toast = useToast()
const confirm = useConfirm()

const valves = ref([])
const loading = ref(true)
const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS }
})
const dialogVisible = ref(false)
const selectedFile = ref(null)

function onFileSelect(event) {
  selectedFile.value = event.files[0]
}

const selectedDatasheetFile = ref(null)

function onDatasheetSelect(event) {
  selectedDatasheetFile.value = event.files[0]
}

const editingId = ref(null)

const emptyForm = {
  manufacturer: '',
  part_number: '',
  component_series: '',
  valve_type: '',
  size_ng: null,
  weight: null,
  rated_flow: null,
  max_flow: null,
  command_value: null,
  command_type: '',
  max_pressure: null
}

const form = reactive({ ...emptyForm })

const errors = reactive({})

function validateForm() {
  Object.keys(errors).forEach(key => delete errors[key])

  if (!form.manufacturer?.trim()) errors.manufacturer = 'Manufacturer is required'
  if (!form.part_number?.trim()) errors.part_number = 'Part Number is required'
  if (!form.component_series?.trim()) errors.component_series = 'Component Series is required'
  if (!form.valve_type?.trim()) errors.valve_type = 'Valve Type is required'
  if (form.size_ng === null || form.size_ng === undefined) errors.size_ng = 'Size NG is required'
  if (form.max_pressure === null || form.max_pressure === undefined) errors.max_pressure = 'Max Pressure is required'

  return Object.keys(errors).length === 0
}

async function fetchValves() {
  try {
    loading.value = true
    const response = await api.get('/valves')
    console.log('DATA VALVE:', response.data)
    valves.value = response.data
  } catch (error) {
    console.error('Failed to fetch valve data:', error)
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  editingId.value = null
  Object.assign(form, emptyForm)
  Object.keys(errors).forEach(key => delete errors[key])
  selectedFile.value = null
  selectedDatasheetFile.value = null
  dialogVisible.value = true
}

function openEditDialog(valve) {
  editingId.value = valve.id
  Object.assign(form, valve) 
  Object.keys(errors).forEach(key => delete errors[key])
  selectedFile.value = null
  selectedDatasheetFile.value = null
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!validateForm()) {
    toast.add({ severity: 'warn', summary: 'the Form is Incomplete', detail: 'Please recheck the required fields', life: 3000 })
    return
  }

  try {
    let valveId = editingId.value

    if (editingId.value) {
      await api.put(`/valves/${editingId.value}`, form)
    } else {
      const response = await api.post('/valves', form)
      valveId = response.data.id
    }

    if (selectedFile.value) {
      await uploadValveImage(valveId, selectedFile.value)
    }

    if (selectedDatasheetFile.value) {           
      await uploadValveDatasheet(valveId, selectedDatasheetFile.value)
    }

    toast.add({ severity: 'success', summary: 'Successful', detail: 'Valve successfully saved', life: 3000 })
    dialogVisible.value = false
    await fetchValves()
  } catch (error) {
    console.error('Failed to save valve:', error)
    toast.add({ severity: 'error', summary: 'Failed', detail: 'An error occurred while saving the data', life: 3000 })
  }
}

function confirmDelete(valve) {
  confirm.require({
    message: `Are you sure you want to delete the "${valve.part_number}" valve ?`,
    header: 'Delete Confirmation',
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: 'Delete',
    rejectLabel: 'Cancel',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.delete(`/valves/${valve.id}`)
        toast.add({ severity: 'success', summary: 'Success', detail: 'Valve was successfully removed', life: 3000 })
        await fetchValves()
      } catch (error) {
        console.error('Failed to delete Valve:', error)
        toast.add({ severity: 'error', summary: 'Failed', detail: 'An error occurred while deleting the data', life: 3000 })
      }
    }
  })
}

onMounted(() => {
  fetchValves()
})
</script>

<template>
  <div class="valve-list-page">
    <Toast />
    <ConfirmDialog />

    <div class="page-header">
      <h2 class="page-title"></h2>
      <Button v-if="authStore.isAdmin" label="Add Valve" icon="pi pi-plus" @click="openCreateDialog" />
    </div>

    <DataTable
      :value="valves"
      :loading="loading"
      stripedRows
      v-model:filters="filters"
      filterDisplay="menu"
      :globalFilterFields="['manufacturer', 'part_number', 'component_series', 'valve_type']"
      paginator
      :rows="5"
      :rowsPerPageOptions="[5, 10, 20, 50]"
      sortMode="multiple"
      scrollable
      scrollHeight="flex"
      class="valve-table"
    >
      <template #header>
        <div class="table-header">
          <IconField>
            <InputIcon class="pi pi-search" />
            <InputText v-model="filters.global.value" placeholder="Search Valve..." />
          </IconField>
        </div>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="manufacturer" header="Manufacturer" sortable />
      <Column field="part_number" header="Part Number" sortable />
      <Column field="component_series" header="Series" sortable />
      <Column field="valve_type" header="Type" sortable />
      <Column field="size_ng" header="Size NG" sortable />
      <Column field="weight" header="Weight" sortable />
      <Column field="rated_flow" header="Rated Flow" sortable />
      <Column field="max_flow" header="Max Flow" sortable />
      <Column field="command_value" header="Command Value" sortable />
      <Column field="command_type" header="Command Type" sortable />
      <Column field="max_pressure" header="Max Pressure" sortable />
      <Column header="Action" :exportable="false" style="min-width: 10rem">
        <template #body="slotProps">
          <Button icon="pi pi-eye" severity="info" text @click="viewDetail(slotProps.data)" />
          <Button icon="pi pi-pencil" severity="secondary" text @click="openEditDialog(slotProps.data)" />
          <Button v-if="authStore.isAdmin" icon="pi pi-trash" severity="danger" text @click="confirmDelete(slotProps.data)" />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editingId ? 'Edit Valve' : 'Tambah Valve'" :modal="true" :style="{ width: '30rem' }">
      <div style="display: flex; flex-direction: column; gap: 1rem;">
        <div>
          <label>Manufacturer</label>
          <InputText v-model="form.manufacturer" :invalid="!!errors.manufacturer" style="width: 100%" />
          <small v-if="errors.manufacturer" style="color: var(--p-red-500)">{{ errors.manufacturer }}</small>
        </div>
        <div>
          <label>Part Number</label>
          <InputText v-model="form.part_number" :invalid="!!errors.part_number" style="width: 100%" />
          <small v-if="errors.part_number" style="color: var(--p-red-500)">{{ errors.part_number }}</small>
        </div>
        <div>
          <label>Component Series</label>
          <InputText v-model="form.component_series" :invalid="!!errors.component_series" style="width: 100%" />
          <small v-if="errors.component_series" style="color: var(--p-red-500)">{{ errors.component_series }}</small>
        </div>
        <div>
          <label>Valve Type</label>
          <InputText v-model="form.valve_type" :invalid="!!errors.valve_type" style="width: 100%" />
          <small v-if="errors.valve_type" style="color: var(--p-red-500)">{{ errors.valve_type }}</small>
        </div>
        <div>
          <label>Size NG</label>
          <InputNumber v-model="form.size_ng" :invalid="!!errors.size_ng" style="width: 100%" />
          <small v-if="errors.size_ng" style="color: var(--p-red-500)">{{ errors.size_ng }}</small>
        </div>
        <div>
          <label>Weight</label>
          <InputNumber v-model="form.weight" :minFractionDigits="0" :maxFractionDigits="2" style="width: 100%" />
        </div>
        <div>
          <label>Rated Flow</label>
          <InputNumber v-model="form.rated_flow" :minFractionDigits="0" :maxFractionDigits="2" style="width: 100%" />
        </div>
        <div>
          <label>Max Flow</label>
          <InputNumber v-model="form.max_flow" :minFractionDigits="0" :maxFractionDigits="2" style="width: 100%" />
        </div>
        <div>
          <label>Command Value</label>
          <InputNumber v-model="form.command_value" :minFractionDigits="0" :maxFractionDigits="2" style="width: 100%" />
        </div>
        <div>
          <label>Command Type</label>
          <InputText v-model="form.command_type" style="width: 100%" />
        </div>
        <div>
          <label>Max Pressure</label>
          <InputNumber v-model="form.max_pressure" :invalid="!!errors.max_pressure" :minFractionDigits="0" :maxFractionDigits="2" style="width: 100%" />
          <small v-if="errors.max_pressure" style="color: var(--p-red-500)">{{ errors.max_pressure }}</small>
        </div>
        <div v-if="authStore.isAdmin">
          <label>Valve Image</label>
          <FileUpload
            mode="basic"
            name="image"
            accept="image/*"
            :maxFileSize="2000000"
            chooseLabel="Choose Image"
            :auto="false"
            customUpload
            @select="onFileSelect"
          />
        </div>
        <div v-if="authStore.isAdmin">
          <label>Datasheet (PDF)</label>
          <FileUpload
            mode="basic"
            name="datasheet"
            accept="application/pdf"
            :maxFileSize="5000000"
            chooseLabel="Choose PDF"
            :auto="false"
            customUpload
            @select="onDatasheetSelect"
          />
        </div>
      </div>

      <template #footer>
        <Button label="Cancel" severity="secondary" @click="dialogVisible = false" />
        <Button label="Save" @click="handleSubmit" />
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.valve-list-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.page-title {
  color: #323638;
  font-size: 2rem;
  font-weight: 700;
  font-family: 'Segoe UI', sans-serif;
  margin: 0;
}

.valve-table {
  flex: 1;
  min-height: 0;
}

.table-header {
  display: flex;
  justify-content: flex-end;
}
</style>