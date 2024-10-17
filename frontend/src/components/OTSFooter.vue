<template>
    <footer class="d-flex align-items-center justify-content-end mt-5 pe-5">
        &copy; OTS - {{ `${new Date().getMonth()} / ${new Date().getFullYear()}` }} 
    </footer>
</template>
<style scoped>
    footer {
        background: black;
        color: white;
        min-height: 3em;
        font-size: 0.8em;
    }
</style>