<script setup lang="ts">

import { ref, computed } from 'vue'
import { getLocalTimeZone, now, ZonedDateTime } from '@internationalized/date'

import { Button } from '@/components/ui/button'
import {
    DialogClose,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog'
import {
    NumberField,
    NumberFieldContent,
    NumberFieldDecrement,
    NumberFieldIncrement,
    NumberFieldInput,
} from '@/components/ui/number-field'

import FieldGroup from '~/components/ui/field/FieldGroup.vue'
import FieldDescription from '~/components/ui/field/FieldDescription.vue'
import Field from '~/components/ui/field/Field.vue'
import FieldSeparator from '~/components/ui/field/FieldSeparator.vue'
import FieldContent from '~/components/ui/field/FieldContent.vue'
import FieldLabel from '~/components/ui/field/FieldLabel.vue'
import FieldSet from '~/components/ui/field/FieldSet.vue'




const maxAmountFreeProducts = 10
const maxHourDuration = 24


const amountDurationHours = ref(1)
const amountFreeProducts = ref(1)
//const amountDurationMinutes = ref(now(getLocalTimeZone()).minute)


const formattedTime = computed((): ZonedDateTime => {
    return now(getLocalTimeZone()).add({ hours: amountDurationHours.value })
})

function handleSubmit() {
  // Hier deine Formular-Logik verarbeiten
}

</script>


<template>


    <DialogContent>

            <form @submit.prevent="handleSubmit">

                <DialogHeader>
                    <DialogTitle>Neues Event erstellen</DialogTitle>
                </DialogHeader>



                <FieldSet>

                    <FieldGroup>
                        <Field>

                            <FieldLabel for="freeProducts">
                                Wähle die Anzahl der Freiprodukte pro Person
                            </FieldLabel>

                            <NumberField id="freeProducts" v-model="amountFreeProducts" :default-value="2" :min="1"
                                :max="maxAmountFreeProducts">
                                <NumberFieldContent>
                                    <NumberFieldDecrement />
                                    <NumberFieldInput />
                                    <NumberFieldIncrement />
                                </NumberFieldContent>
                            </NumberField>
                        </Field>

                        <Field orientation="horizontal">
                            <Button type="button" @click="amountFreeProducts = 2">2 Produkte</Button>
                            <Button type="button" @click="amountFreeProducts = 4">4 Produkte</Button>
                            <Button type="button" @click="amountFreeProducts = 6">6 Produkte</Button>
                            <Button type="button" @click="amountFreeProducts = 10">10 Produkte</Button>
                        </Field>

                        <Field>
                            <FieldContent>
                                <FieldDescription>
                                    Wähle die Menge an Freiprodukte die jeder nutzer erhält. Beim abschluss des Events,
                                    werden
                                    dir alle Produkte in Rechnung gestellt
                                </FieldDescription>
                            </FieldContent>

                        </Field>



                        <FieldSeparator />


                        <Field>
                            <FieldLabel for="eventDuration">
                                Wähle die Laufzeit in Stunden
                            </FieldLabel>

                            <NumberField id="eventDuration" v-model="amountDurationHours" :min="1"
                                :max="maxHourDuration">
                                <NumberFieldContent>
                                    <NumberFieldDecrement />
                                    <NumberFieldInput />
                                    <NumberFieldIncrement />
                                </NumberFieldContent>
                            </NumberField>

                        </Field>

                        <Field orientation="horizontal">
                            <Button type="button" @click="amountDurationHours = 2">2 Stunden</Button>
                            <Button type="button" @click="amountDurationHours = 8">8 Stunden</Button>
                            <Button type="button" @click="amountDurationHours = 16">16 Stunden</Button>
                            <Button type="button" @click="amountDurationHours = 24">24 Stunden</Button>
                        </Field>

                        <Field>
                            <h1>Event Endet am {{ formattedTime.day }}.{{ formattedTime.month }}.{{ formattedTime.year
                            }} um
                                {{ formattedTime.hour.toString().padStart(2, '0') }}:{{
                                    formattedTime.minute.toString().padStart(2, '0') }} Uhr</h1>
                        </Field>


                        <!-- <FieldSeparator /> -->


                    </FieldGroup>
                </FieldSet>

                <DialogFooter>
                    <DialogClose as-child>
                        <Button variant="outline">
                            Abbrechen
                        </Button>
                    </DialogClose>
                    <Button type="submit">
                        Erstellen
                    </Button>
                </DialogFooter>



            </form>

        </DialogContent>


</template>