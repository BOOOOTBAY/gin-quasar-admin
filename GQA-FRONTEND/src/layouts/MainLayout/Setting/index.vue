<template>
    <div>
        <q-btn dense round flat :icon="`${!right ? 'settings' : 'close'}`" @click="right = !right" />
        <q-dialog v-model="right">
            <q-card style="height: 100%; width: 890px; max-width: 80vw;" class="text-center">
                <q-card-section>
                    <div class="text-h6">
                        {{ $t('Customize') + $t('Your') }}
                        {{ gqaFrontend.subTitle }}
                        {{ $t('Display') + $t('Style') }}
                    </div>
                </q-card-section>
                <q-card-section>
                    <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary"
                        align="justify" narrow-indicator>
                        <q-tab name="basicSetting" :label="$t('Basic') + $t('Setting')" />
                        <q-tab name="themeColor" :label="$t('Theme') + $t('Color')" />
                        <q-tab name="themeStyle" :label="$t('Theme') + $t('Style')" alert="negative" />
                    </q-tabs>

                    <q-separator />

                    <q-tab-panels v-model="tab" animated>
                        <q-tab-panel name="basicSetting" class="text-left">
                            <GqaLanguage />
                            <DarkTheme />
                            <SideDrawer />
                        </q-tab-panel>

                        <q-tab-panel name="themeColor">
                            <GqaThemeColor />
                        </q-tab-panel>

                        <q-tab-panel name="themeStyle">
                            <GqaThemeStyle />
                        </q-tab-panel>
                    </q-tab-panels>
                </q-card-section>
            </q-card>
        </q-dialog>
    </div>
</template>

<script setup>
import GqaLanguage from 'src/components/GqaLanguage'
import GqaThemeColor from 'src/components/GqaTheme/GqaThemeColor'
import GqaThemeStyle from 'src/components/GqaTheme/GqaThemeStyle'
import DarkTheme from 'src/components/GqaTheme/DarkTheme.vue';
import SideDrawer from './SideDrawer.vue'
import { ref, computed } from 'vue';
import { useStorageStore } from 'src/stores/storage';

const right = ref(false)
const tab = ref('basicSetting')
const storageStore = useStorageStore();
const gqaFrontend = computed(() => storageStore.GetGqaFrontend())
</script>
