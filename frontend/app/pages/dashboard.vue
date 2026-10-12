<script setup lang="ts">

import { postAuthLogout } from '~/api'
import { useEventUsageStore } from '~/stores/eventStore'




const shoppingCartStore = useShoppingCartStore()
const userStore = useUserStore()
const eventUsageStore = useEventUsageStore()
const productStore = useProductsStore()





userStore.getCurrentUser()
shoppingCartStore.createCurrentShoppingCart()
productStore.fetchProducts()
eventUsageStore.getEventUsageOfCurrentEventAndCurrentUser()



async function cancel() {

    shoppingCartStore.cancelCurrentCart()

    try {
        const response = await postAuthLogout()

    } catch (error) {
        console.error('Error on api auth logout current user', error)
        return null
    }


    await navigateTo('/')
}

async function buy() {

    await shoppingCartStore.checkoutCurrentCart()

    try {
        const response = await postAuthLogout()

    } catch (error) {
        console.error('Error on api auth logout current user', error)
        return null
    }


    await navigateTo('/')

}

</script>


<template>


    <div class="flex flex-col h-dvh p-2 gap-4">

        <Card class="flex-none flex flex-row gap-4 h-22 p-4">

            <Card class="flex-4">
                <NuxtLink to="/user/charge">
                    {{ userStore.currentUser?.name }}
                </NuxtLink>
            </Card>

            <ProductCatalogModal class="flex-2" />

            <Card class="flex-4">
                {{ userStore.currentUser?.balance }}
            </Card>

            <EventModal class="flex-2" />

        </Card>



        <ShoppingCartGrid class="flex-1 min-h-0 overflow-y-auto p-2" />

        <div class="flex-none w-full border-t pt-4">

            <div class="flex items-center justify-between w-full max-w-5xl mx-auto px-4">

                <Button size="lg" class="w-32" @click="cancel">
                    Abbrechen
                </Button>

                <div>
                    <div class="text-2xl font-bold tracking-tight">
                        TOTAL: {{ formatCurrency(shoppingCartStore.currentShoppingCart?.fullPrice) }}
                    </div>
                    <div v-if="shoppingCartStore.currentShoppingCart?.useEvent"
                        class="text-2xl font-bold tracking-tight">
                        Used Free Products: {{ eventUsageStore.eventUsage?.usedFreeAmount }} of
                        {{ eventUsageStore.eventUsage?.availableFreeAmountPerPerson }}
                    </div>
                </div>



                <Button size="lg" class="w-32" @click="buy">
                    Kaufen
                </Button>

            </div>
        </div>


    </div>




</template>
