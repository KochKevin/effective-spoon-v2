import { client } from "~/api/client.gen";


//hey-openapi client settings
export default defineNuxtPlugin( () => {

    client.setConfig({
        
        baseUrl: useRuntimeConfig().public.apiBase,
    })
})