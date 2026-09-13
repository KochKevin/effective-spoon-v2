<script setup lang="ts">


import { Button } from '@/components/ui/button'
import {
    DialogClose,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'


</script>


<template>



    <DialogContent>



        <DialogHeader>
            <DialogTitle>Aktives Event</DialogTitle>
        </DialogHeader>



        <p>Erstellt von Kevin Koch</p>
        <p>Freiprodukte 2 von 5 verfügbar</p>
        <p>Endet am 12.09.2026 18:00 Uhr</p>

        <DialogFooter>
            <DialogClose as-child>
                <Button variant="outline">
                    Schließen
                </Button>
            </DialogClose>
        </DialogFooter>


    </DialogContent>


</template>