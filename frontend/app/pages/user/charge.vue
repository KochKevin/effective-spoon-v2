<script setup lang="ts">
import { Label } from 'reka-ui';
import { postChargementsCurrent } from '~/api';
import Card from '~/components/ui/card/Card.vue';
import NumberField from '~/components/ui/number-field/NumberField.vue';
import NumberFieldContent from '~/components/ui/number-field/NumberFieldContent.vue';
import NumberFieldDecrement from '~/components/ui/number-field/NumberFieldDecrement.vue';
import NumberFieldIncrement from '~/components/ui/number-field/NumberFieldIncrement.vue';
import NumberFieldInput from '~/components/ui/number-field/NumberFieldInput.vue';

const userStore = useUserStore()

const balanceToAdd = ref(15)

async function createStripeCheckoutLink(amount: number): Promise<string | null> {
  try {
    const response = await postChargementsCurrent({
      body: { amountToAdd: amount }
    })

    console.log("paymentLink: ", response.data?.paymentLink ?? null)
    return response.data?.paymentLink ?? null
  } catch (error) {
    console.error('Error post add-balance:', error)
    return null
  }
}



</script>

<template>

    <div>

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

            <Button @click="createStripeCheckoutLink(balanceToAdd)">Jetzt {{ balanceToAdd }}€ mit STRIPE auf dein Profil buchen</Button>

        </Card>

    </div>

</template>