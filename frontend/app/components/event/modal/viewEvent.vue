<script setup lang="ts">


import { Button } from '@/components/ui/button'
import {
    DialogClose,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import { useEventStore } from '~/stores/events';



const eventStore = useEventStore()
const cartStore = useShoppingCartStore()


function availableFreeProducts(): number {
    if (!eventStore.currentEvent) return 0

    return (eventStore.currentEvent.availableAmountPerPerson ?? 0) - (eventStore.currentEvent.eventUsage?.amountUsed ?? 0)
}

function formatEndTime(): string {
    const endTime = eventStore.currentEvent?.eventEndDateTime
    if (!endTime) return 'Kein Event aktiv'

    return endTime
}

    
const useEventSwitch = ref(cartStore.currentShoppingCart?.useEvent) 


function handleUseEventSwitched(checked: boolean) {
    console.log("Changed")
    cartStore.setEventUsage(checked)
}

</script>


<template>

    <div>

        <div class="flex items-center space-x-2">
            <Switch id="use-event" v-model:model-value="useEventSwitch" @update:model-value="handleUseEventSwitched"/>
            <Label for="use-event">Warenkorb nutzt Event</Label>
        </div>

        <p>Erstellt von {{ eventStore.currentEvent?.authorName }}</p>
        <p>Freiprodukte {{ availableFreeProducts() }} von {{ eventStore.currentEvent?.availableAmountPerPerson }}
            verfügbar</p>
        <p>Endet am {{ formatEndTime() }} Uhr</p>

        <DialogFooter>
            <DialogClose as-child>
                <Button variant="outline">
                    Schließen
                </Button>
            </DialogClose>
        </DialogFooter>

    </div>



</template>