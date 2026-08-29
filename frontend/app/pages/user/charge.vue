<script setup lang="ts">
import { postChargementsCurrent, postChargementsCurrentCancel } from '~/api';
import StartChargement from '~/components/chargment/StartChargement.vue';
import ChargementCompleted from '~/components/chargment/ChargementCompleted.vue';
import ChargementQRCode from '~/components/chargment/ChargementQRCode.vue';

const state = ref(0)
const checkoutUrl = ref('')

async function createChargement(amount: number) {
    try {
        const response = await postChargementsCurrent({
            body: { amountToAdd: amount }
        })

        console.log("paymentLink: ", response.data?.paymentLink ?? null);

        checkoutUrl.value = response.data?.paymentLink ?? '';
        state.value = 1;


    } catch (error) {
        console.error('Error on api creating an chargement', error)
        return null
    }
}


async function cancelChargement() {
    
    try {
        const response = await postChargementsCurrentCancel()

    } catch (error) {
        console.error('Error on api canceling current chargement', error)
        return null
    }
    
    console.log("Aufladung abbrechen...")
    await navigateTo('/')
}


const userStore = useUserStore()

async function successfullCompleted() {
    console.log("Aufladung erfolgreich...")
    await userStore.getCurrentUser()
    await navigateTo('/')
}

</script>

<template>


    <StartChargement v-if="state === 0" @create-chargement="createChargement"/>

    <ChargementQRCode v-else-if="state === 1" :qr-value="checkoutUrl" @cancel="cancelChargement" />

    <ChargementCompleted v-else-if="state === 2" @completed="successfullCompleted"/>

    
</template>