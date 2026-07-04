<script setup lang="ts">


defineProps({
    id: String,
    name: String,
    price: Number
})

const shoppingCartStore = useShoppingCartStore()

</script>


<template>

    <Card @click='shoppingCartStore.increaseProduct(shoppingCartStore.shoppingCart?.id ?? "", id ?? "")'>

        <CardHeader>
            <CardTitle>{{ name }}</CardTitle>
        </CardHeader>

        <CardContent>
            <p>{{formatCurrency(price)}}</p>
        </CardContent>

    <!--
        <CardFooter>
            <p>{{id}}</p>
        </CardFooter>
        -->

    </Card>


</template>