import { defineStore } from "pinia";
import { getProducts, type Product } from "~/api";

export const useProductsStore = defineStore('products', {


state: () => {
    return {
        isLoading: true as Boolean,
        products: [] as Product[] | undefined
    }
},


actions: {
    async fetchProducts() {

        this.isLoading = true;

        try {
            const response = await getProducts();

            this.products = response.data
        }
        catch(error) {
            console.error("Error loading products api: ", error);
            
        } finally {
            this.isLoading = false;
        }
        
    }
}


})
