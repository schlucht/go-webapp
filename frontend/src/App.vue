
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

</script>

<style>

</style>