
<script setup lang="ts">
import { Label } from 'reka-ui';
import Card from '~/components/ui/card/Card.vue';
import NumberField from '~/components/ui/number-field/NumberField.vue';
import NumberFieldContent from '~/components/ui/number-field/NumberFieldContent.vue';
import NumberFieldDecrement from '~/components/ui/number-field/NumberFieldDecrement.vue';
import NumberFieldIncrement from '~/components/ui/number-field/NumberFieldIncrement.vue';
import NumberFieldInput from '~/components/ui/number-field/NumberFieldInput.vue';
const userStore = useUserStore()

const balanceToAdd = ref(15)

const emit = defineEmits<{
    createChargement: [amount: number]
}>()

</script>



<template>
    <Card>
            <CardHeader>
                <CardTitle>Lade dein Profil auf {{ userStore.currentUser?.name }}! Du hast {{
                    userStore.currentUser?.balance }}€</CardTitle>
            </CardHeader>


            <NumberField v-model="balanceToAdd" id="balance" :min="0" :format-options="{
                style: 'currency',
                currency: 'EUR',
                currencyDisplay: 'code',
                currencySign: 'accounting',
            }">
                <Label for="balance">Hinzufügen</Label>
                <NumberFieldContent>
                    <NumberFieldDecrement />
                    <NumberFieldInput />
                    <NumberFieldIncrement />
                </NumberFieldContent>
            </NumberField>

            <Button @click="emit('createChargement', balanceToAdd)">Jetzt {{ balanceToAdd }}€ mit STRIPE auf dein Profil
                buchen</Button>


        </Card>
</template>