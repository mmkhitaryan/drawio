<template>
    <div class="brush-settings">
        <v-color-picker canvas-height="70" flat background-color="aliceblue" dot-size="25" hide-inputs @input="setColor" v-model="color"></v-color-picker>
    </div>
</template>

<script>
export default {
    data: () => ({
        current: '#000000'
    }),
    computed: {
      color: {
        get () {
          return this.current;
        },
        set (v) {
          this.current = v
        },
      },
      getColor () {
        if (typeof this.color === 'string') return this.color

        return JSON.stringify(Object.keys(this.color).reduce((color, key) => {
          color[key] = Number(this.color[key].toFixed(2))
          return color
        }, {}), null, 2)
      },
    },
    methods: {
        setColor() {
            this.$emit('change', this.getColor);
        }
    }
}
</script>

<style>
.v-color-picker.v-sheet.theme--light {
    background-color: aliceblue;
}
.brush-settings {
    position: fixed;
    overflow: hidden;
    bottom: 5px;
    left: 5px;

}
</style>