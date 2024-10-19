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
        :value="modelValue"
        :autocomplete="name + '-new'"
        @input="$emit('update:modelValue', $event.target.value)"       
        class="form-control"
    >
  </div>
</template>

<script setup>
    import { defineProps, defineEmits } from 'vue';
    defineProps({
        for:String,
        type:String,
        name:String,
        placeholder:String,
        required:String,
        min:String,
        max:String,
        modelValue:String,
        autocomplete:String,
        label:String,
    });
    /* eslint-disable */
    const emit = defineEmits(["input"]);
    
</script>

<style scoped>

</style>