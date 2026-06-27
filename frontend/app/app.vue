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

        <div class="grid grid-cols-4 gap-4 flex-8">
            <Card>T</Card>
            <Card>T</Card>
            <Card>T</Card>
            <Card>T</Card>
            <Card>T</Card>
      

          </div>


      <AlertDialogFooter class="">
        <AlertDialogCancel>Cancel</AlertDialogCancel>
        <AlertDialogAction>Continue</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>


  </AlertDialog>


</template>
