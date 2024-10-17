<template>
    <div class="container">
        <div class="row">
            <div class="col">
                <h1 class="mt-5">Login</h1>
                <OTSForm  @myevent="submitHandler" name="myform" event="myevent">
                    <OTSTextInput v-model="email" name="email" type="email" placeholder="email" label="Email-Adresse"></OTSTextInput>
                    <OTSTextInput v-model="password" name="password" type="password" placeholder="Passwort" label="Passwort">
                    </OTSTextInput>
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


const password = ref('');
const email = ref('');

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

    const res = await fetch('https://jubilant-space-goldfish-p7wxg999g6vcrgjx-8081.app.github.dev/users/login', body);
    const data = await res.json();

    if(data.error) {
        console.error(data.message);
        return;
    }
    console.log(data.message);
}

</script>

<style scoped></style>