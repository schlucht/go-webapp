<template>
  <div class="mb-3">
    <label :for="name" class="form-label">{{ label }}</label>
    <input
        :type="type"
        :name="name"
        :placeholder="placeholder"
        :required="required"
        :min="min"
        :max="max"
        :value="value"
        :autocomplete="name + '-new'"        
        class="form-control"
    >
  </div>
</template>

<script setup>
    import { defineProps } from 'vue';
    const props = defineProps([
        'for',
        'type',
        'name',
        'placeholder',
        'required',
        'min',
        'max',
        'value',
        'autocomplete',
        'label',
    ]);
    /* eslint-disable */
    props.label = 'Label';
</script>

<style scoped>

</style>