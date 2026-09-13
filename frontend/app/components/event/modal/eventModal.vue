<script setup lang="ts">
import { GiftIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
    Dialog,
    DialogTrigger,
} from '@/components/ui/dialog'

import CreateEvent from './createEvent.vue'
import ViewEvent from './viewEvent.vue'



</script>


<template>
    <Dialog>
        <DialogTrigger as-child>
            <Button variant="outline" size="lg">
                <GiftIcon />
                Event
            </Button>
        </DialogTrigger>

        <ViewEvent />
        <CreateEvent />

    </Dialog>

</template>