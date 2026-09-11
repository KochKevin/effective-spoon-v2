<script setup lang="ts">

import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { postAuthLogout } from '~/api'
import { sseBus } from '~/plugins/02.sse.client'




const shoppingCartStore = useShoppingCartStore()
const userStore = useUserStore()



async function cancel(){
    try {
        const response = await postAuthLogout()

    } catch (error) {
        console.error('Error on api auth logout current user', error)
        return null
    }

    
    await navigateTo('/')
}

async function buy(){
    
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

        </Card>



        <ShoppingCartGrid class="flex-1 min-h-0 overflow-y-auto p-2" />

        <div class="flex-none w-full border-t pt-4">

            <div class="flex items-center justify-between w-full max-w-5xl mx-auto px-4">

                <Button size="lg" class="w-32" @click="cancel">
                    Abbrechen
                </Button>

                <div class="text-2xl font-bold tracking-tight">
                    TOTAL: {{ formatCurrency(shoppingCartStore.currentShoppingCart?.fullPrice) }}
                </div>

                <Button size="lg" class="w-32" @click="buy">
                    Kaufen
                </Button>

            </div>
        </div>


    </div>




</template>
