<script setup lang="ts">

import { sseBus } from './plugins/02.sse.client';


const shoppingCartStore = useShoppingCartStore()


sseBus.on((event) => {
  console.log(event)

  if (event === "user.login") {
    navigateTo("/dashboard")
  }



  if (event === "shoppingcart.update") {
    shoppingCartStore.getCurrentShoppingCart()
  }



})





</script>


<template>

  <div>
    <!-- Markup shared across all pages, ex: NavBar -->
    <NuxtPage />
  </div>


</template>
