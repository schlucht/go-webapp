<template>
  <div class="form-check">
    <input type="checkbox" class="form-check-input"
        :checked="checked"
        :required="required"        
    >
    <label class="form-check-label">{{ label }}</label>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
const props = defineProps([
    'checked',
    'required',
]);
/* eslint-disable */
props.checked = true;
</script>

<style>

</style>