import { defineStore } from "pinia";
import { postShoppingCarts, postShoppingCartsByIdDecrease, postShoppingCartsByIdIncrease, type ShoppingCart } from "~/api";

export const useShoppingCartStore = defineStore('shopping-carts', {


    state: () => {
        return {
            isLoading: true as Boolean,
            shoppingCart: null as ShoppingCart | undefined | null
        }
    },


    actions: {
        async createShoppingCart() {

            this.isLoading = true;

            try {
                const response = await postShoppingCarts();

                this.shoppingCart = response.data
            }
            catch(error) {
                console.error("Error on post on shopping cart api: ", error);

            } finally {
                this.isLoading = false;
            }

        },

        async increaseProduct(cartId : string, productId : string) {

            //this.isLoading = true;

            try {
                const response = await postShoppingCartsByIdIncrease({
                    path: {
                        id: cartId,
                    },

                    query: {
                        productID: productId,
                    }
                });

                this.shoppingCart = response.data
            }
            catch(error) {
                console.error("Error on increase on shopping cart api: ", error);

            } finally {
               // this.isLoading = false;
            }

        },
        async decreaseProduct(cartId : string, productId : string) {

            //this.isLoading = true;

            try {
                const response = await postShoppingCartsByIdDecrease({
                    path: {
                        id: cartId,
                    },

                    query: {
                        productID: productId,
                    }
                });

                this.shoppingCart = response.data
            }
            catch(error) {
                console.error("Error on decrease on shopping cart api: ", error);

            } finally {
               // this.isLoading = false;
            }

        }
    }


})
