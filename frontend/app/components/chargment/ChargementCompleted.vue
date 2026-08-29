<script setup lang="ts">
import { Check } from '@lucide/vue';


const emit = defineEmits<{
    completed: []
}>()

</script>

<template>

    <Card>
        <CardHeader>
            <CardTitle>🎉 Deine Aufladung ist erfolgreich abgeschlossen</CardTitle>
        </CardHeader>

        <CardContent class="flex flex-row justify-center">

            <Check class="text-emerald-500 w-1/3 h-1/3" />

        </CardContent>

        <CardFooter>
            <Button class="w-full" @click="emit('completed')">Abschließen</Button>
        </CardFooter>
    </Card>

</template>