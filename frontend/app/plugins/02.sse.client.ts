import { useEventBus } from "@vueuse/core";
import { getPushes, type SseEvent } from "~/api";


export const sseBus = useEventBus<SseEvent>('sse-event');

export default defineNuxtPlugin(async () => {
    const { stream } = await getPushes();

    try {

        const { stream } = await getPushes();

        (async () => {
            for await (const event of stream) {
                sseBus.emit(event);
            }
        })();

    } catch (err) {
        console.error("Failed to establish SSE connection:", err)
    }
})