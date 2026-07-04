<script setup lang="ts">

import type { LineItem } from '~/api';
import {Icon} from '@iconify/vue'

defineProps<{
  lineItem: LineItem;
}>();

const shoppingCartStore = useShoppingCartStore()

</script>


<template>

   <Item variant="outline">

        <ItemHeader>
            <img
              src="https://picsum.photos/500"
              alt="albert"
              width="128"
              height="128"
              class="aspect-square w-full rounded-sm object-cover"
            >
          </ItemHeader>

        <ItemContent>
            <ItemTitle>Kola</ItemTitle>
            <ItemTitle>{{lineItem.price}}€</ItemTitle>
          
            <ButtonGroup class="">

              <Button variant="outline" size="icon" @click='shoppingCartStore.increaseProduct(shoppingCartStore.shoppingCart?.id ?? "", lineItem.productId ?? "")'>
                <Icon icon="lucide:plus"/>
              </Button>

              <ButtonGroupText>{{ lineItem.amount }}</ButtonGroupText>

               <Button variant="outline" size="icon" @click=''>
                <Icon icon="lucide:minus"/>
              </Button>

            </ButtonGroup>

          </ItemContent>

      </Item>


</template>