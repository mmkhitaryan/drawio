<template>
  <v-app>
    <v-content :style="{'overflowX': 'scroll'}">
      <drawer @draw="progressUp()" ref="drawer"></drawer>
    </v-content>
    <colorpicker @change="setColor"></colorpicker>
    <progress-bar :blocked="blocked" :progress="progress"></progress-bar>
  </v-app>
</template>

<script>
import Colorpicker from './components/Colorpicker.vue';
import ProgressBar from './components/ProgressBar.vue';
import Drawer from './components/Drawer.vue';

export default {
  name: 'App',

  components: { Drawer, Colorpicker, ProgressBar },

  data: () => ({
    progress: 0,
    blocked: false,
    timer: null
  }),
  created() {
    this.timer = setInterval(() => {
      this.progressDown();
    }, 200);
  },
  methods: {
    setColor(color) {
      this.$refs.drawer.setColor(color)
    },
    progressUp() {
      this.progress++;
      if(this.progress >= 240) {
        this.$refs.drawer.blockBrush();
        this.blocked = true;
      }
    },
    progressDown() {
      if(this.progress > 0) {
        this.progress -= 10;
        if(this.progress < 1) {
          this.$refs.drawer.unblockBrush();
          this.blocked = false;
        }
        this.progress = this.progress < 0 ? 0 : this.progress;
      }
    }
  }
};
</script>

<style>
  .v-application--wrap, .v-content, .v-content__wrap{
    max-width: none !important;
  }
</style>
