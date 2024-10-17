<template>
  <div class="mb-3">
    <label :for="name" class="form-label">{{ label }}</label>
    <select v-for="option in items" :value="option.value" :key="option.value">
        <option>
            {{ option.text }}
        </option>
    </select>
  </div>
</template>

<script setup>
    import { defineProps} from 'vue';

    const props = defineProps([
       {items: Object},
       'name',
       'required',
       'label' 
    ]);

    props.label = 'Label';

</script>

<style>

</style>