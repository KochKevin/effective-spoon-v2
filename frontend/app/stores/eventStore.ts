import { defineStore } from "pinia";
import { getEventsCurrent, getEventsCurrentUsageCurrent, postEventsCurrent, type Event, type EventUsage} from "~/api";
import { ZonedDateTime } from '@internationalized/date'


export const useEventUsageStore = defineStore('event-usage', {
 

    state: () => {
        return {
            isLoading: true as Boolean,
            eventUsage: null as EventUsage | undefined | null 
        }
    },


    actions: {


        applyEvent(eventUsage: EventUsage) {
            this.eventUsage = eventUsage
        },
      
        async getEventUsageOfCurrentEventAndCurrentUser() {

            try {
                const response = await getEventsCurrentUsageCurrent();
                
                if (response.response?.status == 200){
                    this.eventUsage = response.data
                }

            }
            catch(error) {
                console.error("Error on get event usage of current event and current user api: ", error);

            } finally {
                this.isLoading = false;
            }

        },

    }


})
