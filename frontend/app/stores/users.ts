import { defineStore } from "pinia";
import { getUsersCurrent, type User } from "~/api";

export const useUserStore = defineStore('users', {


    state: () => {
        return {
            isLoading: true as Boolean,
            currentUser: null as User | undefined | null //Current User is the currently logged in user
        }
    },


    actions: {
        //Current User is the currently logged in user
        async getCurrentUser() {

            this.isLoading = true;

            try {
                const response = await getUsersCurrent();

                this.currentUser = response.data
            }
            catch(error) {
                console.error("Error on get on curent user api: ", error);

            } finally {
                this.isLoading = false;
            }

        },

    }


})
