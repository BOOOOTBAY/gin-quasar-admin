<template>
    <q-dialog v-model="recordDetailVisible" position="top">
        <q-card style="width: 800px; max-width: 80vw;">
            <q-card-section>
                <div class="text-h6">
                    {{ formTypeName }} {{ $t('Api') }}:
                    {{ recordDetail.value.api_path }}
                </div>
            </q-card-section>

            <q-separator />

            <q-card-section>
                <q-form ref="recordDetailForm">
                    <gqa-form-top :recordDetail="recordDetail"></gqa-form-top>
                    <div class="row">
                        <q-input class="col" v-model.number="recordDetail.value.sort" type="number"
                            :rules="[val => val >= 1 || $t('SortRule')]" :label="$t('Sort')" />
                        <q-input class="col" v-model="recordDetail.value.api_group" :label="$t('Api') + $t('Group')"
                            :rules="[val => val && val.length > 0 || $t('NeedInput')]" />
                        <q-select class="col" v-model="recordDetail.value.api_method" :options="methodOptions"
                            :label="$t('Api') + $t('Method')"
                            :rules="[val => val && val.length > 0 || $t('NeedInput')]" />
                        <q-field class="col" dense :label="$t('Status')" stack-label>
                            <template v-slot:control>
                                <q-option-group v-model="recordDetail.value.status" :options="dictOptions.onOff"
                                    color="primary" inline>
                                    <template v-slot:label="opt">
                                        <div class="row items-center">
                                            <span>{{ $t(opt.label) }}</span>
                                        </div>
                                    </template>
                                </q-option-group>
                            </template>
                        </q-field>
                    </div>
                    <div class="row">
                        <q-input class="col" v-model="recordDetail.value.api_path" :label="$t('Api') + $t('Path')"
                            :rules="[val => val && val.length > 0 || $t('NeedInput')]" />
                    </div>
                    <q-input class="col" v-model="recordDetail.value.memo" :label="$t('Memo')" />
                </q-form>
            </q-card-section>

            <q-separator />

            <q-card-actions align="right">
                <q-btn :label="$t('Save')" color="primary" @click="handleAddOrEidt" />
                <q-btn :label="$t('Cancel')" color="negative" v-close-popup />
            </q-card-actions>

            <q-inner-loading :showing="loading">
                <q-spinner-gears size="50px" color="primary" />
            </q-inner-loading>
        </q-card>
    </q-dialog>
</template>


<script setup>
import useRecordDetail from 'src/composables/useRecordDetail'

const emit = defineEmits(['handleFinish'])
const url = {
    add: 'api/add-api',
    edit: 'api/edit-api',
    queryById: 'api/query-api-by-id',
}
const {
    dictOptions,
    formType,
    formTypeName,
    recordDetail,
    recordDetailVisible,
    loading,
    show,
    recordDetailForm,
    handleAddOrEidt
} = useRecordDetail(url, emit)

defineExpose({
    show,
    formType,
    recordDetail
})

const methodOptions = ['POST']
</script>
