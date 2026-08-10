<script setup lang="ts">
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import ItemGroup from '../ui/item/ItemGroup.vue'
import {Icon} from '@iconify/vue'


const shoppingCartStore = useShoppingCartStore()


onMounted( () => {
  

  console.log(shoppingCartStore.isLoading)
  console.log(shoppingCartStore.currentShoppingCart)
})

</script>


<template>


    <p v-if="shoppingCartStore.isLoading">Lade...</p>

    <ItemGroup v-else class="grid grid-cols-4 gap-4 flex-8">


      <ShoppingCartCard v-for="lineItem in shoppingCartStore.currentShoppingCart?.lineItems" :line-item="lineItem"/>



    </ItemGroup>

</template>