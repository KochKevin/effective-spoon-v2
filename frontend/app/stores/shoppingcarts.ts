import { defineStore } from "pinia";
import { getShoppingCartsCurrent, postShoppingCartsCurrent, postShoppingCartsCurrentCheckout, postShoppingCartsCurrentDecrease, postShoppingCartsCurrentIncrease, type ShoppingCart } from "~/api";

export const useShoppingCartStore = defineStore('shopping-carts', {


    state: () => {
        return {
            isLoading: true as Boolean,
            currentShoppingCart: null as ShoppingCart | undefined | null
        }
    },


    actions: {
        async createCurrentShoppingCart() {

            this.isLoading = true;

            try {
                const response = await postShoppingCartsCurrent()

                this.currentShoppingCart = response.data
            }
            catch(error) {
                console.error("Error on post on current shopping cart api: ", error);

            } finally {
                this.isLoading = false;
            }

        },

        async getCurrentShoppingCart() {
            
            this.isLoading = true;

            try {
                const response = await getShoppingCartsCurrent()

                this.currentShoppingCart = response.data
            }
            catch(error) {
                console.error("Error on post on current shopping cart api: ", error);

            } finally {
                this.isLoading = false;
            }
        },

        //Increase Product of current shopping cart
        async increaseProduct(productId : string) {

            //this.isLoading = true;

            try {
                const response = await postShoppingCartsCurrentIncrease({
                    query: {
                        productID: productId,
                    }
                });

                this.currentShoppingCart = response.data
            }
            catch(error) {
                console.error("Error on getting the current shopping cart on its api: ", error);

            } finally {
               // this.isLoading = false;
            }

        },

        //Decrease Product of current shopping cart
        async decreaseProduct(productId : string) {

            //this.isLoading = true;

            try {
                const response = await postShoppingCartsCurrentDecrease({
                    query: {
                        productID: productId,
                    }
                });

                this.currentShoppingCart = response.data
            }
            catch(error) {
                console.error("Error on decrease on current shopping cart api: ", error);

            } finally {
               // this.isLoading = false;
            }

        },

        //Checkout current shopping cart
        async checkoutCurrentCart() {

            //this.isLoading = true;

            try {

                const response = await postShoppingCartsCurrentCheckout();

                this.currentShoppingCart = response.data
            }
            catch(error) {
                console.error("Error on checkout on current shopping cart api: ", error);

            } finally {
               // this.isLoading = false;
            }

        }
    }


})
