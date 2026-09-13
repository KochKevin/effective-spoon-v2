<script setup lang="ts">
import { GiftIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
    Dialog,
    DialogTrigger,
} from '@/components/ui/dialog'

import ViewEvent from './viewEvent.vue'
import { useEventStore } from '~/stores/events.js'
import CreateEvent from './createEvent.vue'
import DialogContent from '~/components/ui/dialog/DialogContent.vue'

const eventStore = useEventStore()




</script>


<template>
    <Dialog>
        <DialogTrigger as-child>
            <Button variant="outline" size="lg">
                <GiftIcon class="text-purple-500 animate-bounce [animation-duration:_1s"v-if="eventStore.currentEvent === null"/>
                <GiftIcon class="text-purple-500 animate-bounce duration-200" v-if="eventStore.currentEvent !== null"/>
                Event
            </Button>
        </DialogTrigger>

        <DialogContent>

            <DialogHeader>
                <DialogTitle v-if="eventStore.currentEvent === null">Neues Event erstellen</DialogTitle>
                <DialogTitle v-if="eventStore.currentEvent !== null">Event</DialogTitle>
            </DialogHeader>

            <CreateEvent v-if="eventStore.currentEvent === null" />
            <ViewEvent v-if="eventStore.currentEvent !== null" />

        </DialogContent>
    </Dialog>

</template>