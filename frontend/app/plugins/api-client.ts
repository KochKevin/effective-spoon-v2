import { client } from "~/api/client.gen";


//hey-openapi client settings
export default defineNuxtPlugin( () => {

    client.setConfig({
        //TODO: Set via .env
        baseUrl: 'https://effective-waddle-4j7xj7vp44pxh7xrp-8080.app.github.dev/'
    })
})