<template>
    <div class="brush-settings">

    </div>
</template>

<script>
export default {
    data: () => ({
        color: "#000000",
    }),
    methods: {
        setColor() {
            this.$emit('change', this.color);
        }
    }
}
</script>