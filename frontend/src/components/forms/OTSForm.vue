<template>
  <form 
    @submit.prevent="submit"
    :ref="name"
    :event="event"
    autocomplete="off" 
    :method="method" 
    :action="action" 
    class="needs-validation" 
    novalidate>
    <slot></slot>
  </form>
</template>

<script setup>
import { defineProps, defineEmits, useTemplateRef } from 'vue';
const props = defineProps([
    "method",
    "action",
    "name",
    'event'
]);
const emit = defineEmits([]);
const itemsRef = useTemplateRef(props['name']);
function submit(){
    let myForm = itemsRef.value;
    if(myForm.checkValidity()) {        
        emit(props['event'], myForm);
    }
    myForm.classList.add('was-validated');
}
</script>

<style>

</style>