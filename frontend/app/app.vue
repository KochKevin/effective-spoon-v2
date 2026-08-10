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
import { useUserStore } from './stores/users';
import { sseBus } from './plugins/02.sse.client';


const shoppingCartStore = useShoppingCartStore()
const userStore = useUserStore()

sseBus.on((event) => {
  console.log(event)

  if (event === "user.login") {
    userStore.getCurrentUser()
    shoppingCartStore.createCurrentShoppingCart()
  }

})





</script>


<template>


  <div class="flex flex-col h-dvh p-2 gap-4">

    <Card class="flex-none flex flex-row gap-4 h-22 p-4">

      <Card class="flex-4">
       {{userStore.currentUser?.name}}
      </Card>

      <ProductCatalogModal class="flex-2"/>

      <Card class="flex-4">
        {{userStore.currentUser?.balance}}
      </Card>

    </Card>

    

    <ShoppingCartGrid class="flex-1 min-h-0 overflow-y-auto p-2" />

    <div class="flex-none w-full border-t pt-4">

      <div class="flex items-center justify-between w-full max-w-5xl mx-auto px-4">

          <Button size="lg" class="w-32">
          Abbrechen
        </Button>

          <div class="text-2xl font-bold tracking-tight">
          TOTAL: {{ formatCurrency(shoppingCartStore.currentShoppingCart?.fullPrice) }}
        </div>

          <Button size="lg" class="w-32" @click="shoppingCartStore.checkoutCurrentCart()">
          Kaufen
        </Button>

        </div>
      </div>


  </div>




</template>
