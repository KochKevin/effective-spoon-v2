import { defineStore } from "pinia";
import { getEventsCurrent, postEventsCurrent, type Event} from "~/api";
import { ZonedDateTime } from '@internationalized/date'


export const useEventStore = defineStore('events', {
 

    state: () => {
        return {
            isLoading: true as Boolean,
            currentEvent: null as Event | undefined | null 
        }
    },


    actions: {
      
        async getCurrentEvent() {

            if (this.currentEvent !== null) {
                return
            }

            this.isLoading = true;

            try {
                const response = await getEventsCurrent();

                this.currentEvent = response.data
            }
            catch(error) {
                console.error("Error on get on curent event api: ", error);

            } finally {
                this.isLoading = false;
            }

        },

        async createCurrentEvent(amountPerPerson: number, endDateTime: ZonedDateTime) {

            this.isLoading = true;

            try {
                const response = await postEventsCurrent({
                    body: {
                        amountPerPerson: amountPerPerson,
                        endDateTime: endDateTime.toString()
                    }
                });

                this.currentEvent = response.data
            }
            catch(error) {
                console.error("Error on create current event api: ", error);

            } finally {
                this.isLoading = false;
            }

        },

    }


})
