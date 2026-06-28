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
import { useProductsStore } from './stores/products';

const productStore = useProductsStore()


onMounted( () => {
  productStore.fetchProducts()

  console.log(productStore.isLoading)
  console.log(productStore.products)
})

</script>


<template>
  <div>
    <NuxtRouteAnnouncer />
    <NuxtWelcome />
  </div>


  <AlertDialog >
    <AlertDialogTrigger>Open</AlertDialogTrigger>

    <AlertDialogContent class="flex flex-col">
      <AlertDialogHeader >
        <AlertDialogTitle>Produkte</AlertDialogTitle>
      </AlertDialogHeader>

       <p v-if="productStore.isLoading">Lade...</p>

        <div v-else class="grid grid-cols-4 gap-4 flex-8">

         
          <Card v-for="product in productStore.products" :key="product.id">
           <p> {{ product.name }} </p>
          </Card>

      

          </div>


      <AlertDialogFooter class="">
        <AlertDialogCancel>Cancel</AlertDialogCancel>
        <AlertDialogAction>Continue</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>


  </AlertDialog>


</template>
