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


const productStore = useProductsStore()


onMounted( () => {
  productStore.fetchProducts()

  console.log(productStore.isLoading)
  console.log(productStore.products)
})

</script>


<template>
 <AlertDialog >
    <AlertDialogTrigger as-child>
        <Button variant="outline">Produkte</Button>
    </AlertDialogTrigger>

    <AlertDialogContent class="flex flex-col">
      <AlertDialogHeader >
        <AlertDialogTitle>Produkte</AlertDialogTitle>
      </AlertDialogHeader>

       <p v-if="productStore.isLoading">Lade...</p>

        <div v-else class="grid grid-cols-4 gap-4 flex-8">

         
          <ProductCatalogCard v-for="product in productStore.products" :key="product.id"
          :id="product.id"
          :name="product.name"
          :price="product.price">
          </ProductCatalogCard>


      

          </div>


      <AlertDialogFooter class="">
        <AlertDialogCancel>Cancel</AlertDialogCancel>
        <AlertDialogAction>Continue</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>


  </AlertDialog>


</template>