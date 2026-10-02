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
import { useEventUsageStore } from './stores/eventStore';


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
