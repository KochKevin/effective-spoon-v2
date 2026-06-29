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


const shoppingCartStore = useShoppingCartStore()


onMounted( () => {
  shoppingCartStore.createShoppingCart()

  console.log(shoppingCartStore.isLoading)
  console.log(shoppingCartStore.shoppingCart)
})

</script>


<template>
 <Card>

    <p v-if="shoppingCartStore.isLoading">Lade...</p>
    <p v-else=> {{shoppingCartStore.shoppingCart}} </p>
 
 
 </Card>


</template>