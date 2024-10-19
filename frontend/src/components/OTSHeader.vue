<template>
  <nav class="navbar bg-dark navbar-expand-lg bg-body-tertiary" data-bs-theme="dark">
  <div class="container-fluid">
    <a class="navbar-brand" href="#">Navbar</a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNav">
      <ul class="navbar-nav me-auto mb-2 mb-lg-0">
        <li class="nav-item">
          <router-link class="nav-link active" aria-current="page" to="/">Home</router-link> 
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="#">Features</router-link>
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="#">Pricing</router-link>
        </li>
        <li class="nav-item">          
          <router-link 
            v-if="store.token ==''"
            class="nav-link" 
            to="/login">Login</router-link>
          <a 
            href="javascript:void(0)"
            v-else
            class="nav-link" 
            @click="logout"
            to="/logout">Logout</a>
        </li>
      </ul>
      <span class="navbar-text">
        {{ store.user.first_name ?? '' }}
      </span>
    </div>
  </div>
</nav>
</template>

<script setup>
  import router from '@/router/router.js';
  import { store } from './store.js';
  import  notie  from 'notie';


  async function logout() {
    const payload = {
      token: store.token,
    };
    const requestOptions = {
      method: "POST",
      body: JSON.stringify(payload),
      headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
    };
    try {
        const res = await fetch('https://jubilant-space-goldfish-p7wxg999g6vcrgjx-8081.app.github.dev/users/logout', requestOptions);
        const data = await res.json();
        if(data.error) {
            console.error(data.message); 
            notie.alert({
                type: 'error',
                text: data.message,
            });           
        } else { 
            notie.alert({
              type: 'success',
              text: data.message,
            });
            store.token = "";
            store.user = {};
            document.cookie = '_site_data=; Path=/;SameSite=Strict;Secure;Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            router.push("/login");
        }
    }catch(err) {

        console.error(err)
    }


    
  }
</script>

<style>

</style>