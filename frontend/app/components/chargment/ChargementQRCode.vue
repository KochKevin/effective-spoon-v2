<script setup lang="ts">
import Card from '~/components/ui/card/Card.vue';
import QrcodeVue from 'qrcode.vue'
import type { Level } from 'qrcode.vue'



const props = defineProps<{
    qrValue: string
}>()

const emit = defineEmits<{
    cancel: []
}>()

const level = ref<Level>('L')
const background = ref(window.getComputedStyle(document.body).getPropertyValue("--primary").trim())
const foreground = ref(window.getComputedStyle(document.body).getPropertyValue("--primary-foreground").trim())


</script>


<template>
    <Card>
        <CardHeader>
            <CardTitle>Scanne den QR Code mit deinem Handy um den Auflade Prozess zu beginn</CardTitle>
        </CardHeader>

        <CardContent class="flex flex-row justify-center">

            <Card class="bg-primary w-1/3 h-1/3">
                <CardContent>
                    <qrcode-vue class="w-full h-full" :value="qrValue" :level="level" :background='background'
                        :foreground='foreground' render-as="svg" />
                </CardContent>
            </Card>

        </CardContent>

        <CardFooter>
            <Button variant="destructive" class="w-full" @click="emit('cancel')">Aufladung Abbrechen</Button>
        </CardFooter>
    </Card>

</template>