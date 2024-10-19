<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-5">Login</h1>
                <OTSForm  @myevent="submitHandler" name="myform" event="myevent">
                    <OTSTextInput 
                        v-model="email" 
                        name="email" 
                        type="email" 
                        placeholder="email" 
                        label="Email-Adresse"
                    ></OTSTextInput>
                    <OTSTextInput 
                        v-model="password" 
                        name="password" 
                        type="password" 
                        placeholder="Passwort" 
                    label="Passwort"
                    ></OTSTextInput>
                    <hr>
                    Email: {{ email }}
                    <hr>
                    <button type="submit" class="btn btn-primary">Submit</button>
                </OTSForm>
            </div>
        </div>
    </div>
</template>

<script setup>
import OTSTextInput from './forms/OTSTextInput.vue';
import OTSForm from './forms/OTSForm.vue';
import { ref } from 'vue';
import { store } from './store.js';
import router from '@/router/router';
import notie from 'notie';


const password = ref('');
const email= ref('');

console.log(password.value, email.value);

async function submitHandler() {

    console.log('submit handler Call');
    const payload = {
        email: email.value,
        password: password.value,
    }
    console.log(JSON.stringify(payload));
    const body = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload),
    }

    try {
        const res = await fetch('https://jubilant-space-goldfish-p7wxg999g6vcrgjx-8081.app.github.dev/users/login', body);
        const data = await res.json();

        if(data.error) {
            console.error(data.message); 
            notie.alert({
                type: 'error',
                text: data.message,
            });           
        } else {            
            store.token = data.data.token.token
            store.user = {
                id: data.data.token.user_id,
                first_name: data.data.token.first_name,
                last_name: data.data.token.last_name,
                email: data.data.token.email,
            };

            let date = new Date();
            let expDay = 1;
            date.setTime(date.getTime() + (expDay * 24 * 60 * 60 * 1000));
            const expires = "expires=" + date.toUTCString();

            document.cookie = "_site_date_"
            + JSON.stringify(date.data)
            + "; "
            + expires
            + "; path=/; SameSite=strict; Secure;";
            router.push("/");
        }
    }catch(err) {
        console.error(err)
    }
    
}

</script>

<style scoped></style>