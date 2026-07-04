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


const shoppingCartStore = useShoppingCartStore()


</script>


<template>


  <div class="flex flex-col h-dvh p-2 gap-4">

    <Card class="flex-none flex flex-row gap-4 h-22 p-4">

      <Card class="flex-4" />

      <ProductCatalogModal class="flex-2"/>

      <Card class="flex-4" />

    </Card>

    

    <ShoppingCartGrid class="flex-1 min-h-0 overflow-y-auto p-2" />

    <div class="flex-none w-full border-t pt-4">

      <div class="flex items-center justify-between w-full max-w-5xl mx-auto px-4">

          <Button size="lg" class="w-32">
          Abbrechen
        </Button>

          <div class="text-2xl font-bold tracking-tight">
          TOTAL: {{ shoppingCartStore.shoppingCart?.fullPrice ?? '0,00' }} €
        </div>

          <Button size="lg" class="w-32">
          Kaufen
        </Button>

        </div>
      </div>


  </div>




</template>
