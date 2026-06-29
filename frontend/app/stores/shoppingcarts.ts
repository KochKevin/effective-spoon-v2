import { defineStore } from "pinia";
import { postShoppingCarts, type ShoppingCart } from "~/api";

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

        }
    }


})
