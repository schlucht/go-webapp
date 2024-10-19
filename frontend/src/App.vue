
<template>
    <OTSHeader /> 
    <div>      
      <router-view></router-view>
    </div>
    <OTSFooter />     
</template>

<script setup>
import OTSHeader from './components/OTSHeader.vue';
import OTSFooter from './components/OTSFooter.vue';
import { onBeforeMount, onMounted } from 'vue';
import { store } from './components/store';

const getCookie = (na) => {
  const c = document.cookie;
  return c.split("; ").reduce((r, v) => {
    const parts = v.split("=");
    return parts[0] === na ? decodeURIComponent(parts[1]) : r;
  }, "")
}

onBeforeMount(() => {
  let data = getCookie("_site_data");
  if (data !== "") {
    let cookieData = JSON.parse(data);
    store.token = cookieData.token.token;
    store.user = {
      id: cookieData.user.id,
      first_name: cookieData.user.first_name,
      last_name: cookieData.user.last_name,
      email: cookieData.user.email,
    }
  }
});

onMounted(async() => {
  const payload = {
    foo: "bar"
  }
   const body = {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'Bearer' + store.token,
        },
        body: JSON.stringify(payload),
    }    
  const res = await fetch('https://jubilant-space-goldfish-p7wxg999g6vcrgjx-8081.app.github.dev/admin/foo', body)
  const data = await res.json();
  console.log(data);
});

</script>

<style>

</style>